package py

import (
	"movieapp/models"
	"strconv"
	"strings"
)

func ConvertToPythonMovie(m models.Movie) (models.PyMovie, error) {
	var movieType string
	var title string
	var releaseYear int
	var dateAdded string
	pyGenre := strings.Join(m.Genre, ", ")
	pyGenre = strings.TrimSpace(pyGenre)
	if m.Name == "" {
		movieType = "Movie"
		title = m.Title
		year, err := strconv.Atoi(m.Release_Date[:4])
		if err != nil {
			return models.PyMovie{}, err
		}
		releaseYear = year
		dateAdded = m.Release_Date

	} else {
		movieType = "TV Show"
		title = m.Name
		year, err := strconv.Atoi(m.First_Air_Date[:4])
		if err != nil {
			return models.PyMovie{}, err
		}
		releaseYear = year

		dateAdded = m.First_Air_Date

	}
	return models.PyMovie{
		Title:       title,
		Type:        &movieType,
		DateAdded:   &dateAdded,
		ReleaseYear: &releaseYear,
		ListedIn:    &pyGenre,
		Description: &m.Overview,
	}, nil
}
