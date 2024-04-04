package main

import (
	"fmt"
	"html/template"
	"io"
	"time"

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

type movie struct {
	Title string
	Slug  string
}

type thread struct {
	Count int
}

type Time struct {
	Time time.Duration
}

func StartServer() {
	e := echo.New()

	// e.Use(middleware.Logger())

	data := map[string]interface{}{}

	e.Renderer = newTemplate()

	// render index.html
	e.GET("/", func(c echo.Context) error {
		data["Movies"] = []movie{
			{Title: "Apocalypse Now", Slug: "apocalypse-now"},
			{Title: "Citizen Kane", Slug: "citizen-kane"},
			{Title: "Parasite", Slug: "parasite-2019"},
		}

		data["Threads"] = []thread{
			{Count: 1},
			{Count: 2},
			{Count: 4},
			{Count: 8},
			{Count: 10},
			{Count: 12},
			{Count: 16},
		}
		return c.Render(200, "index.html", data)
	})

	e.POST("/search", func(c echo.Context) error {
		// TODO: introduce some rate limiting

		chosen := c.FormValue("movie")
		maxUsers := c.FormValue("users")
		threads := c.FormValue("threads")

		fmt.Println("Starting search for", chosen, "with", maxUsers, "users and", threads, "threads")

		movies, time := Scraper(chosen, maxUsers, threads)

		fmt.Println(time)

		// map the first 5 movies to movie structs
		data["ReturnedMovies"] = []movie{}

		// right now it shows the movie you chose, but we should remove it
		for i := 0; i < min(5, len(movies)); i++ {
			data["ReturnedMovies"] = append(data["ReturnedMovies"].([]movie), movie{Title: movies[i], Slug: movies[i]})
		}

		data["Time"] = Time{Time: time}

		// ideally here we render what movies we found
		return c.Render(200, "results.html", data)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
