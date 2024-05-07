package endpoints

import (
	"letterboxd-rec/services"
	"letterboxd-rec/templates/components"
	"letterboxd-rec/templates/partials"
	"letterboxd-rec/utils"
	"net/http"
	"strconv"
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

	mux.HandleFunc("POST /open-dropdown", func(w http.ResponseWriter, r *http.Request) {
		componet := components.OpenDropdown()
		componet.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /close-dropdown", func(w http.ResponseWriter, r *http.Request) {
		componet := components.CloseDropdown()
		componet.Render(r.Context(), w)
	})

	mux.HandleFunc("PUT /update-settings", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		maxUsers, err := strconv.Atoi(r.FormValue("maxUsers"))
		userErr := ""

		if err != nil || maxUsers < 5 || maxUsers > 5000 {
			userErr = "invalid users"
		} else {
			utils.MaxUsers = maxUsers
		}

		threads, err := strconv.Atoi(r.FormValue("threads"))
		threadsErr := ""

		if err != nil || threads < 0 || threads > 32 {
			threadsErr = "invalid threads"
		} else {
			utils.Threads = threads
		}

		// close settings if no errors
		if userErr == "" && threadsErr == "" {
			w.Header().Set("HX-Retarget", "#dropdown")
			w.Header().Set("HX-Reswap", "outerHTML swap:500ms")

			components.CloseDropdown().Render(r.Context(), w)
		}
		components.ErrorLabel(userErr, threadsErr).Render(r.Context(), w)
	})
}
