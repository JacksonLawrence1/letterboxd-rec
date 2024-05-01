package main

import (
	"net/http"

	"letterboxd-rec/handlers"
)

// command to run the server
// templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."

func main() {
	h := handlers.New()

	http.ListenAndServe(":8080", h)
}
