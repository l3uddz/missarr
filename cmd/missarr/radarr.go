package main

import (
	"database/sql"
	"fmt"
	"github.com/l3uddz/missarr/migrate"
	"github.com/l3uddz/missarr/radarr"
	"github.com/rs/zerolog/log"
	"time"
)

type RadarrCmd struct {
	Limit           int           `default:"10" help:"How many items to search for before stopping"`
	LastSearched    time.Duration `default:"672h" help:"How long before an item can be searched again"`
	LastReleaseDate time.Duration `default:"72h" help:"How long before an item can be considered missing based on release date"`
	SkipRefresh     bool          `default:"false" help:"Retrieve current missing from radarr"`
	Delay           time.Duration `default:"0s" help:"Delay between search requests"`
}

func (r *RadarrCmd) Run(c *config, db *sql.DB, mg *migrate.Migrator) error {
	// validate flags
	if r.Limit == 0 {
		r.Limit = 10
	}

	// init
	sc, err := radarr.New(&c.Radarr, db, mg)
	if err != nil {
		return fmt.Errorf("initialise radarr: %w", err)
	}

	// retrieve movies
	var rm []radarr.MovieResponse
	var fm []radarr.Movie
	var us, rs int

	if !r.SkipRefresh {
		rm, err = sc.Movies()
		if err != nil {
			return fmt.Errorf("retrieving movies: %w", err)
		}
		log.Info().
			Int("size", len(rm)).
			Msg("Retrieved movies")

		// refresh datastore
		us, rs, fm, err = sc.RefreshStore(rm, time.Now().Add(-r.LastReleaseDate))
		if err != nil {
			return fmt.Errorf("missing to store: %w", err)
		}
		log.Info().
			Int("missing_movies", us).
			Int("completed_movies", rs).
			Msg("Refreshed datastore")
	} else {
		fm, err = sc.GetAll()
		if err != nil {
			return fmt.Errorf("get all from datastore: %w", err)
		}
	}

	// search for movies
	eligibleSearchDate := time.Now().Add(-r.LastSearched)
	for _, s := range fm {
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
			Int("movie", s.Id).
			Time("release_date", s.ReleaseDate).
			Msg("Searching...")

		searchId, err := sc.Search(&s)
		if err != nil {
			return fmt.Errorf("search movie: %w", err)
		}

		log.Info().
			Int("movie", s.Id).
			Time("release_date", s.ReleaseDate).
			Int("search_id", searchId).
			Msg("Search requested")

		// update store
		now := time.Now()
		s.SearchDate = &now
		if err := sc.UpdateStore([]radarr.Movie{s}); err != nil {
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
