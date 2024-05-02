package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"letterboxd-rec/utils"
)

type MovieSearchData struct {
	Results []utils.Movie
}

// Get a single movie's information from TMDB by its id
func GetMovieInfo(id int) (utils.Movie, error) {
	movie := utils.Movie{}

	// endpoint
	url := "https://api.themoviedb.org/3/movie/" + fmt.Sprint(id) + "?language=en-US"
	res, err := request_TMDB(url)

	if err != nil {
		return movie, err
	}

	// filter out the data we need from the json response
	defer res.Body.Close()

	// decode the json
	body, _ := io.ReadAll(res.Body)

	// only need the title, id and poster path (for now)
	json.Unmarshal(body, &movie)

	// if it unmarshals correctly, return the movie, otherwise return an error
	if movie.Id == 0 || movie.Title == "" || movie.Poster_path == "" || movie.Release_date == "" {
		return movie, fmt.Errorf("movie id not found on tmdb")
	}

	// get the year from the release date
	movie.Release_date = movie.Release_date[:4]

	return movie, nil
}

// Get a list of movies from the search term using TMDB's search endpoint
func SearchForMovies(term string) ([]utils.Movie, error) {
	searchResults := MovieSearchData{}

	url := "https://api.themoviedb.org/3/search/movie?query=" + strings.ReplaceAll(term, " ", "%20") + "&include_adult=false&language=en-US&page=1"
	res, err := request_TMDB(url)

	if err != nil {
		return searchResults.Results, err
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return searchResults.Results, err
	}

	json.Unmarshal(body, &searchResults)

	return searchResults.Results, nil
}
