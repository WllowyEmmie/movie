package py

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"movieapp/models"
	"net/http"
	"time"
)

func RecommendFromFavorites(favorites []models.Movie) ([]string, error) {
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
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("python API error: %s", string(body))
	}
	defer resp.Body.Close()
	var result []string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}
