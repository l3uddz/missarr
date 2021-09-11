package radarr

import (
	"fmt"
	"time"
)

func (c *Client) UpdateStore(movie []Movie) error {
	return c.store.Upsert(movie)
}

func (c *Client) GetAll() ([]Movie, error) {
	return c.store.GetAll()
}

func (c *Client) RefreshStore(data []MovieResponse, maxReleaseDate time.Time) (int, int, []Movie, error) {
	// filter movies
	movies := make([]Movie, 0)
	sm := make(map[int]int)

	for _, m := range data {
		// skip if movie matches conditions
		switch {
		case !m.Monitored:
			continue
		case !m.IsAvailable:
			continue
		case m.HasFile, m.SizeOnDisk > 0:
			continue
		}

		// determine release date to use
		releaseDate := m.Added
		switch {
		case !m.DigitalRelease.IsZero():
			releaseDate = m.DigitalRelease
		case !m.InCinemas.IsZero():
			releaseDate = m.InCinemas
		}

		// skip movie if the release date is not before the max release date
		if releaseDate.After(maxReleaseDate) {
			continue
		}

		// add movie
		sm[m.Id] = 1
		movies = append(movies, Movie{
			Id:          m.Id,
			ReleaseDate: releaseDate,
			SearchDate:  nil,
		})
	}

	moviesSize := len(movies)

	// store movies in datastore
	if err := c.store.Upsert(movies); err != nil {
		return 0, 0, nil, fmt.Errorf("upsert: %w", err)
	}

	// retrieve movies from datastore
	em, err := c.store.GetAll()
	if err != nil {
		return moviesSize, 0, movies, fmt.Errorf("get all: %w", err)
	}

	// generate movies to remove
	moviesToRemove := make([]Movie, 0)
	finalMovies := make([]Movie, 0)
	for _, m := range em {
		if _, ok := sm[m.Id]; !ok {
			moviesToRemove = append(moviesToRemove, m)
			continue
		}
		finalMovies = append(finalMovies, m)
	}

	// remove movies from datastore
	moviesToRemoveSize := len(moviesToRemove)
	if moviesToRemoveSize > 0 {
		if err := c.store.Delete(moviesToRemove); err != nil {
			return moviesSize, 0, movies, fmt.Errorf("remove no longer missing: %w", err)
		}
	}

	return moviesSize, moviesToRemoveSize, finalMovies, nil
}
