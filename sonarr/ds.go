package sonarr

import (
	"fmt"
	"time"
)

func (c *Client) MissingToStore(episodes []Episode) (int, error) {
	// sort episodes into series
	sm := make(map[string]time.Time)
	seasons := make([]Series, 0)

	for _, e := range episodes {
		// skip if episode is not monitored, or we already have a file
		if !e.Monitored || e.HasFile {
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

	// store seasons in datastore
	if err := c.store.Upsert(seasons); err != nil {
		return 0, fmt.Errorf("upsert: %w", err)
	}

	return len(sm), nil
}
