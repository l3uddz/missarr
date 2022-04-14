package sonarr

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/missarr/util"
	"github.com/lucperkins/rek"
	"net/url"
	"time"
)

type systemStatus struct {
	Version string
}

func (c *Client) getSystemStatus() (*systemStatus, error) {
	// send request
	resp, err := rek.Get(util.JoinURL(c.apiURL, "system", "status"), rek.Client(c.http), rek.Headers(c.apiHeaders))
	if err != nil {
		return nil, fmt.Errorf("request system status: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate system status response: %s", resp.Status())
	}

	// decode response
	b := new(systemStatus)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode system status response: %w", err)
	}

	return b, nil
}

type Episode struct {
	SeriesId                 int       `json:"seriesId"`
	EpisodeFileId            int       `json:"episodeFileId"`
	SeasonNumber             int       `json:"seasonNumber"`
	EpisodeNumber            int       `json:"episodeNumber"`
	Title                    string    `json:"title"`
	AirDate                  string    `json:"airDate"`
	AirDateUtc               time.Time `json:"airDateUtc"`
	HasFile                  bool      `json:"hasFile"`
	Monitored                bool      `json:"monitored"`
	UnverifiedSceneNumbering bool      `json:"unverifiedSceneNumbering"`
	Id                       int       `json:"id"`
}

func (c *Client) Missing(cutoff bool) ([]Episode, error) {
	// prepare request
	missingType := "missing"
	if cutoff {
		missingType = "cutoff"
	}
	reqURL, err := util.URLWithQuery(util.JoinURL(c.apiURL, "wanted", missingType), url.Values{
		"page":     []string{"1"},
		"pageSize": []string{"100000"},
		"sortDir":  []string{"desc"},
		"sortKey":  []string{"airDateUtc"},
	})
	if err != nil {
		return nil, fmt.Errorf("create missing url: %w", err)
	}

	// send request
	resp, err := rek.Get(reqURL, rek.Client(c.http), rek.Headers(c.apiHeaders))
	if err != nil {
		return nil, fmt.Errorf("request missing: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate missing response: %s", resp.Status())
	}

	// decode response
	b := new(struct {
		TotalRecords int       `json:"totalRecords"`
		Records      []Episode `json:"records"`
	})
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode missing response: %w", err)
	}

	return b.Records, nil
}

func (c *Client) Search(series *Series) (int, error) {
	// prepare request
	p := new(struct {
		Name         string `json:"name"`
		SeriesId     int    `json:"seriesId"`
		SeasonNumber int    `json:"seasonNumber"`
	})

	p.Name = "SeasonSearch"
	p.SeriesId = series.Id
	p.SeasonNumber = series.Season

	// send request
	resp, err := rek.Post(util.JoinURL(c.apiURL, "command"), rek.Client(c.http), rek.Headers(c.apiHeaders),
		rek.Json(p))
	if err != nil {
		return 0, fmt.Errorf("request search: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 201 {
		return 0, fmt.Errorf("validate search response: %s", resp.Status())
	}

	// decode response
	b := new(struct {
		Id int `json:"id"`
	})
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return 0, fmt.Errorf("decode search response: %w", err)
	}

	return b.Id, nil
}
