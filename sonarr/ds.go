package sonarr

import (
	"fmt"
	"time"
)

func (c *Client) UpdateStore(season []Series) error {
	return c.store.Upsert(season)
}

func (c *Client) GetAll() ([]Series, error) {
	return c.store.GetAll()
}

func (c *Client) RefreshStore(episodes []Episode, allowSpecials bool, maxAirDate time.Time) (int, int, []Series, error) {
	// sort episodes into series
	sm := make(map[string]time.Time)
	seasons := make([]Series, 0)

	for _, e := range episodes {
		// skip if episode is not monitored, or we already have a file
		if !e.Monitored || e.HasFile {
			continue
		}

		// skip specials
		if !allowSpecials && e.SeasonNumber == 0 {
			continue
		}

		// skip episode if the air date is not before the max air date
		if e.AirDateUtc.After(maxAirDate) {
			continue
		}

		// series season before?
		k := fmt.Sprintf("%v_%v", e.SeriesId, e.SeasonNumber)
		if _, ok := sm[k]; ok {
			continue
		}
		// add to map
		sm[k] = e.AirDateUtc

		seasons = append(seasons, Series{
			Id:         e.SeriesId,
			Season:     e.SeasonNumber,
			AirDate:    e.AirDateUtc,
			SearchDate: nil,
		})
	}

	seasonsSize := len(seasons)

	// store seasons in datastore
	if err := c.store.Upsert(seasons); err != nil {
		return 0, 0, nil, fmt.Errorf("upsert: %w", err)
	}

	// retrieve seasons from datastore
	es, err := c.store.GetAll()
	if err != nil {
		return seasonsSize, 0, seasons, fmt.Errorf("get all: %w", err)
	}

	// generate seasons to remove
	seasonsToRemove := make([]Series, 0)
	finalSeasons := make([]Series, 0)
	for _, s := range es {
		k := fmt.Sprintf("%v_%v", s.Id, s.Season)
		if _, ok := sm[k]; !ok {
			seasonsToRemove = append(seasonsToRemove, s)
		}
		finalSeasons = append(finalSeasons, s)
	}

	// remove seasons from datastore
	seasonsToRemoveSize := len(seasonsToRemove)
	if seasonsToRemoveSize > 0 {
		if err := c.store.Delete(seasonsToRemove); err != nil {
			return seasonsSize, 0, seasons, fmt.Errorf("remove no longer missing: %w", err)
		}
	}

	return seasonsSize, seasonsToRemoveSize, finalSeasons, nil
}
