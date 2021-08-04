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
	epps, err := sc.Missing()
	if err != nil {
		return fmt.Errorf("retrieving missing: %w", err)
	}

	log.Info().Int("size", len(epps)).Msg("Retrieved missing")

	// store missing in datastore
	if err := sc.MissingToStore(epps); err != nil {
		return fmt.Errorf("missing to store: %w", err)
	}

	return nil
}
