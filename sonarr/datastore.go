package sonarr

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/l3uddz/missarr/migrate"
	"time"

	// sqlite3 driver
	_ "modernc.org/sqlite"
)

type datastore struct {
	*sql.DB
}

var (
	//go:embed migrations
	migrations embed.FS
)

func newDatastore(db *sql.DB, mg *migrate.Migrator) (*datastore, error) {
	// migrations
	if err := mg.Migrate(&migrations, "sonarr"); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &datastore{db}, nil
}

const sqlUpsert = `
INSERT INTO series (series, season, air_date, search_date)
VALUES (?, ?, ?, ?)
ON CONFLICT (series) DO UPDATE SET
	air_date = excluded.air_date
    , season = excluded.season
    , search_date = COALESCE(series.search_date, excluded.search_date)
`

func (store *datastore) upsert(tx *sql.Tx, series int, season int, airDate time.Time, searchDate *time.Time) error {
	_, err := tx.Exec(sqlUpsert, series, season, airDate, searchDate)
	return err
}

type Series struct {
	Id         int
	Season     int
	AirDate    time.Time
	SearchDate *time.Time
}

func (store *datastore) Upsert(series []Series) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	for _, s := range series {
		if err = store.upsert(tx, s.Id, s.Season, s.AirDate, s.SearchDate); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				panic(rollbackErr)
			}

			return err
		}
	}

	return tx.Commit()
}
