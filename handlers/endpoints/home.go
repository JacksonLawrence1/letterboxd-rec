package endpoints

import (
	"letterboxd-rec/services"
	"letterboxd-rec/templates/pages"
	"letterboxd-rec/templates/partials"
	"letterboxd-rec/utils"
	"net/http"
)

func addToTrending(trendingMovies *[]utils.Movie) error {
	if len(*trendingMovies) == 0 {
		trending, err := services.GetTrending()

		if err != nil {
			return err
		}

		*trendingMovies = append(*trendingMovies, trending...)
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

		partials.Homepage(trendingMovies).Render(r.Context(), w)
	})
}
