package handlers

import (
	"letterboxd-rec/handlers/endpoints"
	"net/http"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Index and Home route
	endpoints.HomeHandler(mux)

	// Search handler
	endpoints.SearchHandler(mux)

	// Recommend handler
	endpoints.RecommendHandler(mux)

	return mux
}
