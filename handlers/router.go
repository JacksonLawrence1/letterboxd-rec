package handlers

import (
	"letterboxd-rec/templates/pages"
	"letterboxd-rec/templates/partials"
	"net/http"
)

func New(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index := pages.Index()
		index.Render(r.Context(), w)
	})

	// Home
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		partials.SearchBar().Render(r.Context(), w)
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		pages.About().Render(r.Context(), w)
	})

	// Search handler
	SearchHandler(mux)

	// Recommend handler
	RecommendHandler(mux)
}
