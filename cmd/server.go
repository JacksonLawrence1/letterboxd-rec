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
		// movie := c.FormValue("movie")

		movie := "parasite-2019"

		// Use the Scraper
		movies, _ := Scraper(movie, 5, 4) // movie, maxusers, threads

		for i := 0; i < len(movies); i++ {
			data["Movies"] = append(data["Movies"].([]Movie), Movie{Title: movies[i], Slug: movies[i]})
		}

		return c.Render(200, "recommendations.html", nil)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
