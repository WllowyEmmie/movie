package server

import (
	"movieapp/models"
	"movieapp/tmdb"
	"net/http"
	"sort"
	"sync"

	"github.com/labstack/echo/v4"
)

func SetUpRoutes(e *echo.Echo, api *tmdb.TMDbClient) {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World!")
	})
	e.GET("api/movies/search", func(c echo.Context) error {

		query := c.QueryParam("q")

		if query == "" {
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"error": "Missing Query Parameter"}, " ")
		}

		results, err := api.SearchMovies(query)
		if err != nil {
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": err.Error()}, " ")
		}

		return c.JSONPretty(http.StatusOK, results, " ")
	})
	e.GET("api/tv/search", func(c echo.Context) error {
		query := c.QueryParam("q")

		if query == "" {
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"error": "Missing Query Parameters"}, " ")
		}
		results, err := api.SearchSeries(query)
		if err != nil {
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": err.Error()}, " ")
		}
		return c.JSONPretty(http.StatusOK, results, " ")
	})
	e.GET("api/person/search", func(c echo.Context) error {
		query := c.QueryParam("q")

		if query == "" {
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"error": "Missing Query Parameters"}, " ")
		}
		results, err := api.SearchPerson(query)
		if err != nil {
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": err.Error()}, " ")
		}
		return c.JSONPretty(http.StatusOK, results, " ")
	})
	e.GET("api/search", func(c echo.Context) error {
		query := c.QueryParam("q")

		if query == "" {
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"error": "Missing Query Parameter"}, " ")
		}
		var wg sync.WaitGroup
		var movies, series []models.Movie

		var movieErr, seriesErr error

		wg.Add(2)

		go func() {
			defer wg.Done()
			movies, movieErr = api.SearchMovies(query)
		}()
		go func() {
			defer wg.Done()
			series, seriesErr = api.SearchSeries(query)
		}()
		wg.Wait()
		if movieErr != nil {
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": movieErr.Error()}, " ")
		}

		if seriesErr != nil {
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"error": seriesErr.Error()}, " ")
		}

		results := append(movies, series...)
		sort.Slice(results, func(i, j int) bool {
			return results[i].VoteAverage > results[j].VoteAverage
		})

		return c.JSONPretty(http.StatusOK, results, " ")
	})
}
