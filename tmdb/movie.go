package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"movieapp/models"
	"net/http"
	"net/url"
)

type TMDbClient struct {
	BaseURL string
	Token   string
	Client  *http.Client
}

func (t *TMDbClient) Get(path string, v interface{}) error {
	url := fmt.Sprintf("%s%s", t.BaseURL, path)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+t.Token)
	req.Header.Add("Accept", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("TMDb API error: %s", string(body))
	}

	return json.NewDecoder(resp.Body).Decode(v)

}
func (t *TMDbClient) SearchMovies(query string) ([]models.Movie, error) {
	type response struct {
		Results []models.Movie `json:"results"`
	}

	var res response
	encodedQuery := url.QueryEscape(query)
	err := t.Get(fmt.Sprintf("/search/movie?query=%s", encodedQuery), &res)
	if err != nil {
		return nil, err
	}

	for i := range res.Results {
		if res.Results[i].Poster_Path != "" {
			res.Results[i].Type, res.Results[i].Poster_Path = designateTypeAndImage(res.Results[i].Name, res.Results[i].Poster_Path)
		}
		isTV := res.Results[i].Name != ""
		for _, id := range res.Results[i].Genre_Ids {
			if name := GetGenreName(id, isTV); name != "" {
				res.Results[i].Genre = append(res.Results[i].Genre, name)
			}

		}
	}
	return res.Results, nil

}
func (t *TMDbClient) SearchSeries(query string) ([]models.Movie, error) {
	type response struct {
		Results []models.Movie `json:"results"`
	}

	var res response
	encodedQuery := url.QueryEscape(query)
	err := t.Get(fmt.Sprintf("/search/tv?query=%s", encodedQuery), &res)

	if err != nil {
		return nil, err
	}
	for i := range res.Results {
		if res.Results[i].Poster_Path != "" {
			res.Results[i].Type, res.Results[i].Poster_Path = designateTypeAndImage(res.Results[i].Name, res.Results[i].Poster_Path)

		}
		isTV := res.Results[i].Name != ""
		for _, id := range res.Results[i].Genre_Ids {
			if name := GetGenreName(id, isTV); name != "" {
				res.Results[i].Genre = append(res.Results[i].Genre, name)
			}

		}
	}
	return res.Results, nil
}
func (t *TMDbClient) SearchPerson(query string) ([]models.Actor, error) {
	type response struct {
		Results []models.Actor `json:"results"`
	}
	var res response
	encodedQuery := url.QueryEscape(query)
	err := t.Get(fmt.Sprintf("/search/person?query=%s", encodedQuery), &res)
	if err != nil {
		return nil, err
	}
	for i := range res.Results {
		for j := range res.Results[i].Known_for {
			res.Results[i].Known_for[j].Type, res.Results[i].Known_for[j].Poster_Path = designateTypeAndImage(res.Results[i].Known_for[j].Name, res.Results[i].Known_for[j].Poster_Path)

			isTV := res.Results[i].Known_for[j].Name != ""
			for _, id := range res.Results[i].Known_for[j].Genre_Ids {
				if name := GetGenreName(id, isTV); name != "" {
					res.Results[i].Known_for[j].Genre = append(res.Results[i].Known_for[j].Genre, name)
				}

			}
		}

	}
	return res.Results, nil
}
func designateTypeAndImage(name, poster_path string) (mediaType, imageUrl string) {
	const baseImageUrl = "https://image.tmdb.org/t/p/w500"
	if name != "" {
		mediaType = "Serial"
	} else {
		mediaType = "Movie"
	}
	if poster_path != "" {
		imageUrl = fmt.Sprintf("%s%s", baseImageUrl, poster_path)
	}
	return
}
