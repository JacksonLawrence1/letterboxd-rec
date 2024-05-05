package tmdb

import (
	"encoding/json"
	"fmt"
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
	body, err := requestTMDBBody(url)

	if err != nil {
		return movie, err
	}

	// only need the title, id and poster path (for now)
	json.Unmarshal(body, &movie)

	// if it unmarshals correctly, return the movie, otherwise return an error
	if movie.Id == 0 || movie.Title == "" || movie.Poster_path == "" || movie.Release_date == "" {
		return movie, fmt.Errorf("movie id not found on tmdb")
	}

	return movie, nil
}

func GetTrendingMovies() ([]utils.Movie, error) {
	trendingMovies := MovieSearchData{}

	url := "https://api.themoviedb.org/3/movie/popular?language=en-US&page=1"
	body, err := requestTMDBBody(url)

	if err != nil {
		return trendingMovies.Results, err
	}

	json.Unmarshal(body, &trendingMovies)
	return trendingMovies.Results, nil
}

// Get a list of movies from the search term using TMDB's search endpoint
func SearchForMovies(term string) ([]utils.Movie, error) {
	searchResults := MovieSearchData{}

	url := "https://api.themoviedb.org/3/search/movie?query=" + strings.ReplaceAll(term, " ", "%20") + "&include_adult=false&language=en-US&page=1"
	body, err := requestTMDBBody(url)

	if err != nil {
		return searchResults.Results, err
	}

	json.Unmarshal(body, &searchResults)
	return searchResults.Results, nil
}
