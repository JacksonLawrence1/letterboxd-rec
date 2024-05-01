package handlers

import (
	"letterboxd-rec/components"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	// Serve static assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index := components.Index()
		index.Render(r.Context(), w)
	})

	// Home
	mux.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		components.Home().Render(r.Context(), w)
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		components.About().Render(r.Context(), w)
	})

	// Search handler
	SearchHandler(mux)

	// Recommend handler
	RecommendHandler(mux)

	return mux
}
