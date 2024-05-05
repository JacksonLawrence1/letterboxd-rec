package handlers

import (
	"letterboxd-rec/templates/pages"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Index and Home route
	HomeHandler(mux)

	// About route
	mux.HandleFunc("POST /about", func(w http.ResponseWriter, r *http.Request) {
		pages.About().Render(r.Context(), w)
	})

	// Search handler
	SearchHandler(mux)

	// Recommend handler
	RecommendHandler(mux)

	return mux
}
