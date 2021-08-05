package main

import (
	"database/sql"
	"fmt"
	"github.com/l3uddz/missarr/migrate"
	"github.com/l3uddz/missarr/sonarr"
	"github.com/rs/zerolog/log"
)

type SonarrCmd struct {
	Limit int `default:"10" help:"How many items to search for before stopping"`
}

func (r *SonarrCmd) Run(c *config, db *sql.DB, mg *migrate.Migrator) error {
	// validate flags
	if r.Limit == 0 {
		r.Limit = 10
	}

	// init
	sc, err := sonarr.New(&c.Sonarr, db, mg)
	if err != nil {
		return fmt.Errorf("initialise sonarr: %w", err)
	}

	// retrieve missing
	se, err := sc.Missing()
	if err != nil {
		return fmt.Errorf("retrieving missing: %w", err)
	}
	log.Info().
		Int("size", len(se)).
		Msg("Retrieved missing episodes")

	// refresh datastore
	us, rs, fs, err := sc.RefreshStore(se)
	if err != nil {
		return fmt.Errorf("missing to store: %w", err)
	}
	log.Info().
		Int("incomplete_seasons", us).
		Int("completed_seasons", rs).
		Msg("Refreshed datastore")

	// search for series
	for _, s := range fs {
		// limit reached
		if r.Limit == 0 {
			break
		}

		// search for season
		log.Info().
			Int("series", s.Id).
			Int("season", s.Season).
			Time("air_date", s.AirDate).
			Msg("Searching...")

		// decrease limit
		r.Limit--
	}

	return nil
}
