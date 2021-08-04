package sonarr

import (
	"fmt"
	"time"
)

func (c *Client) MissingToStore(episodes []Episode) error {
	// sort episodes into series
	sm := make(map[int]time.Time)
	series := make([]Series, 0)
	for _, e := range episodes {
		// seen before?
		if _, ok := sm[e.SeriesId]; ok {
			continue
		}
		// add to map
		sm[e.SeriesId] = e.AirDateUtc

		series = append(series, Series{
			Id:         e.SeriesId,
			AirDate:    e.AirDateUtc,
			SearchDate: nil,
		})
	}

	c.log.Debug().
		Int("series", len(sm)).
		Msg("Sorted missing into unique series")

	// store episodes in datastore
	if err := c.store.Upsert(series); err != nil {
		return fmt.Errorf("upsert: %w", err)
	}

	return nil
}
