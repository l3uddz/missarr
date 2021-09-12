package radarr

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
	if err := mg.Migrate(&migrations, "radarr"); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return &datastore{db}, nil
}

const sqlUpsert = `
INSERT INTO movies (movie, release_date, search_date)
VALUES (?, ?, ?)
ON CONFLICT (movie) DO UPDATE SET
	release_date = excluded.release_date
    , search_date = CASE
        WHEN excluded.release_date > movies.release_date THEN NULL
        ELSE COALESCE(excluded.search_date, movies.search_date)
		END
`

func (store *datastore) upsert(tx *sql.Tx, movie int, releaseDate time.Time, searchDate *time.Time) error {
	_, err := tx.Exec(sqlUpsert, movie, releaseDate, searchDate)
	return err
}

type Movie struct {
	Id          int
	ReleaseDate time.Time
	SearchDate  *time.Time
}

func (store *datastore) Upsert(movies []Movie) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	for _, m := range movies {
		if err = store.upsert(tx, m.Id, m.ReleaseDate, m.SearchDate); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				panic(rollbackErr)
			}

			return err
		}
	}

	return tx.Commit()
}

const sqlGetAll = `
SELECT movie, release_date, search_date
FROM movies
ORDER BY release_date DESC
`

func (store *datastore) GetAll() (movies []Movie, err error) {
	rows, err := store.Query(sqlGetAll)
	if err != nil {
		return movies, err
	}

	defer rows.Close()
	for rows.Next() {
		movie := Movie{}
		err = rows.Scan(&movie.Id, &movie.ReleaseDate, &movie.SearchDate)
		if err != nil {
			return movies, err
		}

		movies = append(movies, movie)
	}

	return movies, rows.Err()
}

const sqlDelete = `
DELETE FROM movies WHERE movie = ?
`

func (store *datastore) delete(tx *sql.Tx, movie Movie) error {
	_, err := tx.Exec(sqlDelete, movie.Id)
	if err != nil {
		return fmt.Errorf("delete: %w", err)
	}
	return nil
}

func (store *datastore) Delete(movies []Movie) error {
	tx, err := store.Begin()
	if err != nil {
		return err
	}

	for _, m := range movies {
		if err = store.delete(tx, m); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				panic(rollbackErr)
			}

			return err
		}
	}

	return tx.Commit()
}
