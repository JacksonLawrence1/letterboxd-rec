package handlers

import (
	"letterboxd-rec/components"
	"letterboxd-rec/services"
	"letterboxd-rec/utils"
	"net/http"
	"strconv"
)

var searchResultsMap = map[int]*utils.Movie{}

func SearchHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		searchTerm := r.FormValue("movie")

		// results of the search on TMDB
		searchResults, err := services.Search(searchTerm)

		if err != nil {
			components.ErrorPage(404, err.Error()).Render(r.Context(), w)
			return
		}

		// map the tmdb id to the movie so we can use it when recommending
		for _, movie := range searchResults {
			searchResultsMap[movie.Id] = &movie
		}

		// probably should filter results not found on letterboxd first
		component := components.Results(searchResults)
		component.Render(r.Context(), w)
	})
}

func RecommendHandler(mux *http.ServeMux) {
	movies := []utils.Movie{}
	movie := utils.Movie{}

	mux.HandleFunc("POST /recommend", func(w http.ResponseWriter, r *http.Request) {
		tmdbId, err := strconv.Atoi(r.FormValue("tmdb-id"))

		if err != nil || searchResultsMap[tmdbId] == nil {
			components.ErrorPage(404, "Error while searching, please try again.").Render(r.Context(), w)
			return
		}

		movie = *searchResultsMap[tmdbId]

		movies, err = services.Recommend(movie)

		if err != nil {
			components.ErrorPage(404, err.Error()).Render(r.Context(), w)
			return
		}

		recommendationPanel := components.RecommendationPanel(&movie, movies, len(movies) < utils.ItemsToShow)
		recommendationPanel.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /loadMore", func(w http.ResponseWriter, r *http.Request) {
		TMDBMovieInfo, isFull := services.GetMoreRecommendationsInfo()

		if len(TMDBMovieInfo) > 0 {
			movies = append(movies, TMDBMovieInfo...)
		}

		updatedRecommendations := components.Recommendations(movies, isFull)
		updatedRecommendations.Render(r.Context(), w)
	})
}
