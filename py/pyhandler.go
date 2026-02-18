package py

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"movieapp/models"
	"movieapp/tmdb"
	"net/http"
	"strings"
	"time"
)

func RecommendFromFavorites(favorites []models.Movie) ([]models.Movie, error) {
	url := "http://127.0.0.1:8000/api/recommend_multiple"
	var pyFavorites []models.PyMovie
	for _, movie := range favorites {
		pyMovie, err := ConvertToPythonMovie(movie)
		if err != nil {
			return nil, err
		}
		pyFavorites = append(pyFavorites, pyMovie)
	}

	payload, err := json.Marshal(pyFavorites)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 50 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("python API error: %s", string(body))
	}

	var result []models.PyMovie
	var tmdbresult []models.Movie
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	for _, movie := range result {
		title := movie.Title
		if *movie.Type == "movie" {
			moviesearch, err := tmdb.SetUpApi().SearchMovies(title)
			if err != nil {
				return nil, err
			}
			for _, exact := range moviesearch {
				if strings.EqualFold(exact.Title, movie.Title) {
					tmdbresult = append(tmdbresult, exact)
				}
			}
		} else {
			moviesearch, err := tmdb.SetUpApi().SearchSeries(title)
			if err != nil {
				return nil, err
			}
			for _, exact := range moviesearch {
				if strings.EqualFold(exact.Name, movie.Title) {
					tmdbresult = append(tmdbresult, exact)
				}
			}
		}
	}

	return tmdbresult, nil
}
