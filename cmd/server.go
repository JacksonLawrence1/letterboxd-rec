package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Data struct {
	Title string
}

func StartServer() {
	h1 := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("./views/index.html"))
		templ.Execute(w, nil)
	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")

		// search for film title here

		htmlStr := fmt.Sprintf("<h1>Search results for: %s</h1>", title)
		templ, _ := template.New("search").Parse(htmlStr)

		templ.Execute(w, nil)
	}

	http.HandleFunc("/", h1)
	http.HandleFunc("/search", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
