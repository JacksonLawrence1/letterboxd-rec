package endpoints

import (
	"letterboxd-rec/services"
	"letterboxd-rec/templates/partials"
	"net/http"
)

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

		// filter results not on letterboxd
		component := partials.Results("search results", searchResults)
		component.Render(r.Context(), w)
	})
}
