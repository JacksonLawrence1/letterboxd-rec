package services

import (
	"fmt"
	"letterboxd-rec/services/letterboxd"
	"letterboxd-rec/services/tmdb"
	"letterboxd-rec/utils"
)

var movieData = utils.MovieData{}

func Search(term string) ([]utils.Movie, error) {
	// get movie data
	searchResults, err := tmdb.SearchForMovies(term)

	if err != nil {
		return searchResults, err
	}

	// don't show movies that don't have enough fans on letterboxd
	filteredResults := letterboxd.FilterLetterboxdMoviesByFans(searchResults, utils.MaxUsers)

	if len(filteredResults) == 0 {
		return filteredResults, fmt.Errorf("no results found on letterboxd")
	}

	return filteredResults, nil
}

func Recommend(movie utils.Movie) ([]utils.Movie, error) {
	// Reset the movie data
	movieData = utils.MovieData{Movie: movie, Pointer: 0, Slugs: []string{}}

	// this users the movie's film slug, make sure you look up the correct one
	users := letterboxd.ScrapeUsers(&movieData.Movie)

	// Scrapes each part of the user's profile to get their 4 favourites
	// This part takes the longest time as we need to scrape a new page for each user
	movieFrequencyMap := letterboxd.ScrapeFavourites(users)

	if len(movieFrequencyMap) == 0 {
		return []utils.Movie{}, fmt.Errorf("no movies found")
	}

	movieData.Slugs = utils.SortByFrequency(movieData.Movie.Slug, movieFrequencyMap)

	// first batch should not be full
	batch, _ := GetMoreRecommendations()

	// get the first batch of movies
	return batch, nil
}

func GetMoreRecommendations() ([]utils.Movie, bool) {
	// get the relevant ids by scraping the letterboxd page
	tmdbIds := letterboxd.ConvertMovieSlugs(movieData.Slugs[movieData.Pointer:movieData.Increment()])

	movieBatch := []utils.Movie{}

	// get the movie info from TMDB using the ids
	for _, tmdbId := range tmdbIds {
		movie, err := tmdb.GetMovieInfo(tmdbId)

		// if there's an error, skip the movie
		if err == nil {
			movieBatch = append(movieBatch, movie)
		}
	}

	return movieBatch, movieData.IsFull()
}
