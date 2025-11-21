package tmdb

import (
	"net/http"
	"os"
)

func SetUpApi() *TMDbClient {
	token := os.Getenv("API_BEARER_TOKEN")
	if token == "" {
		panic("Bearer token not set")
	}
	url := "https://api.themoviedb.org/3"

	return &TMDbClient{
		BaseURL: url,
		Token:   token,
		Client:  &http.Client{},
	}

}
