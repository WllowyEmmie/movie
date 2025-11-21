package tmdb

import (
	"fmt"
	"log"
	"sync"
)

var (
	MovieGenreMap = map[int]string{}
	TvGenreMap    = map[int]string{}
	mu            sync.RWMutex
)

func (t *TMDbClient) LoadGenres() error {
	if err := t.loadGenreList("/genre/movie/list", MovieGenreMap); err != nil {
		return fmt.Errorf("failed to load movie genres: %w", err)
	}
	if err := t.loadGenreList("/genre/tv/list", TvGenreMap); err != nil {
		return fmt.Errorf("failed to load tv genres: %w", err)
	}
	log.Println("✅ Genres successfully loaded.")
	return nil
}

func (t *TMDbClient) loadGenreList(endpoint string, targetMap map[int]string) error {
	type genreResponse struct {
		Genres []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"genres"`
	}
	var res genreResponse
	if err := t.Get(endpoint, &res); err != nil {
		return fmt.Errorf("failed to load genre: %w", err)
	}
	mu.Lock()
	defer mu.Unlock()
	for _, g := range res.Genres {
		targetMap[g.ID] = g.Name
	}
	return nil
}
func GetGenreName(id int, isTV bool) string {
	mu.RLock()
	defer mu.RUnlock()
	if isTV {
		return TvGenreMap[id]
	}
	return MovieGenreMap[id]
}
