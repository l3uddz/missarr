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
ON CONFLICT (series, season) DO UPDATE SET
	air_date = excluded.air_date
    , search_date = CASE
        WHEN excluded.air_date > series.air_date THEN NULL
        ELSE COALESCE(excluded.search_date, series.search_date)
		END
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

const sqlGetAll = `
SELECT series, season, air_date, search_date
FROM series
ORDER BY air_date DESC
`

func (store *datastore) GetAll() (seasons []Series, err error) {
	rows, err := store.Query(sqlGetAll)
	if err != nil {
		return seasons, err
	}

	defer rows.Close()
	for rows.Next() {
		series := Series{}
		err = rows.Scan(&series.Id, &series.Season, &series.AirDate, &series.SearchDate)
		if err != nil {
			return seasons, err
		}

		seasons = append(seasons, series)
	}

	return seasons, rows.Err()
}

const sqlDelete = `
DELETE FROM series WHERE series = ? AND season = ?
`

func (store *datastore) delete(tx *sql.Tx, series Series) error {
	_, err := tx.Exec(sqlDelete, series.Id, series.Season)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (store *datastore) Delete(series []Series) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	for _, s := range series {
		if err = store.delete(tx, s); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				panic(rollbackErr)
			}

			return err
		}
	}

	return tx.Commit()
}
