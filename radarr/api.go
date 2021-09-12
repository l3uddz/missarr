package radarr

import (
	"encoding/json"
	"fmt"
	"github.com/l3uddz/missarr/util"
	"github.com/lucperkins/rek"
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

type MovieResponse struct {
	SizeOnDisk     int       `json:"sizeOnDisk"`
	Status         string    `json:"status"`
	InCinemas      time.Time `json:"inCinemas"`
	DigitalRelease time.Time `json:"digitalRelease"`
	Year           int       `json:"year"`
	HasFile        bool      `json:"hasFile"`
	Monitored      bool      `json:"monitored"`
	IsAvailable    bool      `json:"isAvailable"`
	Added          time.Time `json:"added"`
	Id             int       `json:"id"`
}

func (c *Client) Movies() ([]MovieResponse, error) {
	// send request
	resp, err := rek.Get(util.JoinURL(c.apiURL, "movie"), rek.Client(c.http), rek.Headers(c.apiHeaders))
	if err != nil {
		return nil, fmt.Errorf("request movies: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate movies response: %s", resp.Status())
	}

	// decode response
	b := new([]MovieResponse)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode movies response: %w", err)
	}

	return *b, nil
}

func (c *Client) Search(movie *Movie) (int, error) {
	// prepare request
	p := new(struct {
		Name     string `json:"name"`
		MovieIds []int  `json:"movieIds"`
	})

	p.Name = "MoviesSearch"
	p.MovieIds = []int{movie.Id}

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
