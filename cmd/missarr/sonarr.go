package main

import (
	"fmt"
	"github.com/l3uddz/missarr/sonarr"
	"github.com/rs/zerolog/log"
	"time"
)

type SonarrCmd struct {
	Limit int `default:"10" help:"How many items to search for before stopping"`
}

func (r *SonarrCmd) Run(c *config) error {
	// validate flags
	if r.Limit == 0 {
		r.Limit = 10
	}

	// init
	sc, err := sonarr.New(&c.Sonarr)
	if err != nil {
		return fmt.Errorf("initialise sonarr: %w", err)
	}

	// retrieve missing
	epps, err := sc.Missing()
	if err != nil {
		return fmt.Errorf("retrieving missing: %w", err)
	}

	log.Info().Int("size", len(epps)).Msg("Retrieved missing")

	// sort missing into seasons
	series := make(map[int]time.Time)
	for _, e := range epps {
		// seen before?
		if _, ok := series[e.SeriesId]; ok {
			continue
		}
		// add to map
		series[e.SeriesId] = e.AirDateUtc
	}
	log.Info().Int("series", len(series)).Msg("Sorted missing into unique series")

	return nil
}
