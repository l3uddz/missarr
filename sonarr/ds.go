package sonarr

import (
	"fmt"
	"time"
)

func (c *Client) MissingToStore(episodes []Episode) (int, error) {
	// sort episodes into series
	sm := make(map[int]time.Time)
	series := make([]Series, 0)

	for _, e := range episodes {
		// skip if episode is not monitored, or we already have a file
		if !e.Monitored || e.HasFile {
			continue
		}

		// seen before?
		if _, ok := sm[e.SeriesId]; ok {
			continue
		}
		// add to map
		sm[e.SeriesId] = e.AirDateUtc

		series = append(series, Series{
			Id:         e.SeriesId,
			Season:     e.SeasonNumber,
			AirDate:    e.AirDateUtc,
			SearchDate: nil,
		})
	}

	// store episodes in datastore
	if err := c.store.Upsert(series); err != nil {
		return 0, fmt.Errorf("upsert: %w", err)
	}

	return len(sm), nil
}
