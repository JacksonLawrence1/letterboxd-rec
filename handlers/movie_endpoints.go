package handlers

import (
	"letterboxd-rec/components"
	"letterboxd-rec/services"
	"letterboxd-rec/utils"
	"net/http"
)

func SearchHandler(mux *http.ServeMux) {
	mux.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		// the movie the user is searching for
		movie := r.FormValue("movie")

		// results of the search on TMDB
		searchResults := services.Search(movie)

		component := components.Results(searchResults)
		component.Render(r.Context(), w)
	})
}

func RecommendHandler(mux *http.ServeMux) {
	moviesData := utils.MovieData{Movie: utils.Movie{}, Pointer: 0, Slugs: []string{}}
	movies := []utils.Movie{}

	mux.HandleFunc("POST /recommend", func(w http.ResponseWriter, r *http.Request) {
		movieSlug := r.FormValue("movie")

		// get movie data
		movieData, err := services.LookUpMovie(movieSlug)

		if err != nil {
			components.ErrorPage(404, "Movie not found on TMDB").Render(r.Context(), w)
		}

		// Use the Scraper to get the movie slugs
		movieSlugs, err := services.Scraper(movieSlug) // movie, maxusers, threads

		if err != nil {
			components.ErrorPage(404, "Movie not found on Letterboxd").Render(r.Context(), w)
		}

		moviesData = utils.MovieData{Movie: movieData, Pointer: 0, Slugs: movieSlugs}

		// Use the LookUp to get the movie info
		// this also overwrites the movie data
		movies, _ = services.LookUpMovies(&moviesData)

		if len(movies) == 0 {
			components.ErrorPage(404, "No movies found on Letterboxd").Render(r.Context(), w)
		}

		component := components.RecommendationPanel(&moviesData, movies)
		component.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /loadMore", func(w http.ResponseWriter, r *http.Request) {
		TMDBMovieInfo, _ := services.LookUpMovies(&moviesData)

		// add the data to our existing movies
		movies = append(movies, TMDBMovieInfo...)
		components.Recommendations(movies, moviesData.IsFull()).Render(r.Context(), w)
	})
}
