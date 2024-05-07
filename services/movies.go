package services

import (
	"fmt"
	"letterboxd-rec/services/letterboxd"
	"letterboxd-rec/services/tmdb"
	"letterboxd-rec/utils"
)

var movieData = utils.RecommendationData{}

func Search(term string) ([]utils.Movie, error) {
	// get movie data
	searchResults, err := tmdb.SearchForMovies(term)

	if err != nil {
		return searchResults, err
	}

	// don't show movies that don't have enough fans on letterboxd
	filteredResults := letterboxd.FilterLetterboxdMoviesByFans(searchResults, utils.MaxUsers, utils.Threads)

	if len(filteredResults) == 0 {
		return filteredResults, fmt.Errorf("no results found on letterboxd")
	}

	return filteredResults, nil
}

func GetTrending() ([]utils.Movie, error) {
	// get trending movies
	trendingMovies, err := tmdb.GetTrendingMovies()

	if err != nil {
		return trendingMovies, err
	}

	filteredResults := letterboxd.FilterLetterboxdMoviesByFans(trendingMovies, utils.MaxUsers, utils.Threads)

	if len(filteredResults) == 0 {
		return filteredResults, fmt.Errorf("no trending movies found on letterboxd")
	}

	filteredResults = utils.SortByFans(filteredResults)

	// only get the top 8 trending movies (or less if we don't have enough)
	return filteredResults[:min(4, len(filteredResults))], nil
}

func Recommend(movie utils.Movie, maxUsers int, threads int) ([]utils.Movie, bool, error) {
	// Reset the movie data
	movieData = utils.RecommendationData{Movie: movie, Slugs: []string{}}

	// Scrapes the users who have the movie in their favourites
	users := letterboxd.ScrapeUsers(&movieData.Movie, maxUsers, threads)

	utils.Progress.Message = "Getting user favourites"

	// Scrapes each part of the user's profile to get their 4 favourites
	// This part takes the longest time as we need to scrape a new page for each user
	movieFrequencyMap := letterboxd.ScrapeFavourites(users, threads)

	if len(movieFrequencyMap) == 0 {
		return []utils.Movie{}, true, fmt.Errorf("no movies found")
	}

	utils.Progress.Message = "finalising recommendations"

	movieData.Slugs = utils.SortByFrequency(movieData.Movie.Slug, movieFrequencyMap)

	// first batch should not be full
	batch, _, isFull := GetMoreRecommendations(0, threads)

	// get the first batch of movies
	return batch, isFull, nil
}

func GetMoreRecommendations(start int, threads int) ([]utils.Movie, int, bool) {
	end := min(start+utils.ItemsToShow, len(movieData.Slugs))

	// get the relevant ids by scraping the letterboxd page
	tmdbIds := letterboxd.ConvertSlugToTMDBId(movieData.Slugs[start:end], utils.Threads)

	movieBatch := []utils.Movie{}

	// get the movie info from TMDB using the ids
	for i, id := range tmdbIds {
		// this means it wasn't found on letterboxd
		if id == 0 {
			continue
		}

		movie, err := tmdb.GetMovieInfo(id)
		movie.Slug = movieData.Slugs[start+i]

		// if there's an error, skip the movie
		if err == nil {
			movieBatch = append(movieBatch, movie)
		}
	}

	return movieBatch, end + 1, end >= len(movieData.Slugs)
}
