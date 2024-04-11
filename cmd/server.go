package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
}

func StartServer() {
	e := echo.New()

	// e.Use(middleware.Logger())

	e.Renderer = newTemplate()

	data := map[string]interface{}{}

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index.html", nil)
	})

	e.POST("/search", func(c echo.Context) error {
		// first we need to sanitise the input

		// possibilities
		// - only have a small selection of movies
		// - use TMBD API to search for the movies (might be difficult right now)
		movie := c.FormValue("movie")

		// Use the Scraper
		movies, err := Scraper(movie, 5, 4) // movie, maxusers, threads

		if err != nil {
			data["Error"] = "Movie not found"
			return c.Render(200, "error.html", data)
		}

		// Use the LookUp to get the movie info
		movieData, _ := LookUpMovies(movies)

		data["Movies"] = movieData

		return c.Render(200, "recommendations.html", data)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
