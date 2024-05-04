package main

import (
	"net/http"

	"letterboxd-rec/handlers"
)

// command to run the server
// templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

func main() {
	mux := http.NewServeMux()

	// Serve static assets
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	handlers.New(mux)

	http.ListenAndServe(":8080", mux)
}
