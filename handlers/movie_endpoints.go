package handlers

import (
	"letterboxd-rec/services"
	"letterboxd-rec/templates/pages"
	"letterboxd-rec/templates/partials"
	"letterboxd-rec/utils"
	"net/http"
	"strconv"
)

var searchResultsMap = map[int]*utils.Movie{}

func addToTrending(trendingMovies *[]utils.Movie) error {
	if len(*trendingMovies) == 0 {
		trending, err := services.GetTrending()

		if err != nil {
			return err
		}

		for _, movie := range trending {
			searchResultsMap[movie.Id] = &movie
			*trendingMovies = append(*trendingMovies, movie)
		}
	}
	return nil
}

func HomeHandler(mux *http.ServeMux) {
	// cache the trending movies
	trendingMovies := []utils.Movie{}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		addToTrending(&trendingMovies)

		pages.Index(trendingMovies).Render(r.Context(), w)
	})

	// Home
	mux.HandleFunc("POST /home", func(w http.ResponseWriter, r *http.Request) {
		err := addToTrending(&trendingMovies)

		if err != nil {
			pages.ErrorPage(500, "Error while fetching movies.").Render(r.Context(), w)
			return
		}

		for _, movie := range trendingMovies {
			searchResultsMap[movie.Id] = &movie
		}

		partials.Homepage(trendingMovies).Render(r.Context(), w)
	})
}

func SearchHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		searchTerm := r.FormValue("movie")

		// results of the search on TMDB
		searchResults, err := services.Search(searchTerm)

		if err != nil {
			w.Header().Set("HX-Retarget", "#validation")
			partials.ValidationError("No Movies Found").Render(r.Context(), w)
			return
		}

		// map the tmdb id to the movie so we can use it when recommending
		for _, movie := range searchResults {
			searchResultsMap[movie.Id] = &movie
		}

		// probably should filter results not found on letterboxd first
		component := partials.Results("search results", searchResults)
		component.Render(r.Context(), w)
	})
}

func RecommendHandler(mux *http.ServeMux) {
	movies := []utils.Movie{}
	movie := utils.Movie{}

	mux.HandleFunc("POST /recommend", func(w http.ResponseWriter, r *http.Request) {
		tmdbId, err := strconv.Atoi(r.FormValue("tmdb-id"))

		if err != nil || searchResultsMap[tmdbId] == nil {
			pages.ErrorPage(500, "Error while searching, please try again.").Render(r.Context(), w)
			return
		}

		// Cache the selected movie
		movie = *searchResultsMap[tmdbId]

		// Get the recommendations for the selected movie
		movies, err = services.Recommend(movie)

		if err != nil {
			pages.ErrorPage(500, err.Error()).Render(r.Context(), w)
			return
		}

		// Render recommendation panel
		recommendationPanel := partials.RecommendationPanel(&movie, movies)
		recommendationPanel.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /isFull", func(w http.ResponseWriter, r *http.Request) {
		partials.LoadMore().Render(r.Context(), w)
	})

	mux.HandleFunc("POST /loadMore", func(w http.ResponseWriter, r *http.Request) {
		TMDBMovieInfo, isFull := services.GetMoreRecommendations()

		if isFull {
			w.Header().Set("HX-Trigger", "moreResults")
		}

		// Add the new recommendations
		if len(TMDBMovieInfo) > 0 {
			updatedRecommendations := partials.Recommendations(TMDBMovieInfo)
			updatedRecommendations.Render(r.Context(), w)
		}
	})
}
