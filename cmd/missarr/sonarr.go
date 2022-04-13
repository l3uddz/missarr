package main

import (
	"database/sql"
	"fmt"
	"github.com/l3uddz/missarr/migrate"
	"github.com/l3uddz/missarr/sonarr"
	"github.com/rs/zerolog/log"
	"time"
)

type SonarrCmd struct {
	Limit        int           `default:"10" help:"How many items to search for before stopping"`
	LastSearched time.Duration `default:"672h" help:"How long before an item can be searched again"`
	LastAirDate  time.Duration `default:"72h" help:"How long before an item can be considered missing based on air date"`
	AllowSpecial bool          `default:"false" help:"Allow specials to be considered missing"`
	SkipRefresh  bool          `default:"false" help:"Retrieve current missing from sonarr"`
	Delay        time.Duration `default:"0s" help:"Delay between search requests"`
	Cutoff       bool          `default:"false" help:"Search Cutoff Unmet Items"`
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
	var se []sonarr.Episode
	var fs []sonarr.Series
	var us, rs int

	if !r.SkipRefresh {
		se, err = sc.Missing(r.Cutoff)
		if err != nil {
			return fmt.Errorf("retrieving missing: %w", err)
		}
		log.Info().
			Int("size", len(se)).
			Msg("Retrieved missing episodes")

		// refresh datastore
		us, rs, fs, err = sc.RefreshStore(se, r.AllowSpecial, time.Now().Add(-r.LastAirDate), r.Cutoff)
		if err != nil {
			return fmt.Errorf("missing to store: %w", err)
		}
		log.Info().
			Int("incomplete_seasons", us).
			Int("completed_seasons", rs).
			Msg("Refreshed datastore")
	} else {
		fs, err = sc.GetAll()
		if err != nil {
			return fmt.Errorf("get all from datastore: %w", err)
		}
	}

	// search for series
	eligibleSearchDate := time.Now().Add(-r.LastSearched)
	for _, s := range fs {
		// search date is eligible ?
		if s.SearchDate != nil && s.SearchDate.After(eligibleSearchDate) {
			continue
		}

		// limit reached
		if r.Limit == 0 {
			break
		}

		// search for season
		log.Debug().
			Int("series", s.Id).
			Int("season", s.Season).
			Time("air_date", s.AirDate).
			Msg("Searching...")

		searchId, err := sc.Search(&s)
		if err != nil {
			return fmt.Errorf("search series: %w", err)
		}

		log.Info().
			Int("series", s.Id).
			Int("season", s.Season).
			Time("air_date", s.AirDate).
			Int("search_id", searchId).
			Msg("Search requested")

		// update store
		now := time.Now()
		s.SearchDate = &now
		if err := sc.UpdateStore([]sonarr.Series{s}); err != nil {
			return fmt.Errorf("update store: %w", err)
		}

		// sleep for delay
		if r.Limit > 0 {
			time.Sleep(r.Delay)
		}

		// decrease limit
		r.Limit--
	}

	return nil
}
