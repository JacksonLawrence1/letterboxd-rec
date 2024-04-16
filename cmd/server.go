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

	e.Static("/images", "images")
	e.Static("/css", "css")

	data := map[string]interface{}{}

	// essentially the "backend" data of the movies
	movieData := MovieData{"", 0, []int{}}
	data["Movie"] = &movieData.movie

	// the frontend data of the movies, e.g. full title, release date and poster
	movies := []Movie{}
	data["Movies"] = &movies

	// whether there are more movies to load
	data["HasMore"] = true

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", nil)
	})

	e.POST("/search", func(c echo.Context) error {
		data["HasMore"] = true

		// possibilities
		// - only have a small selection of movies
		// - use TMBD API to search for the movies (might be difficult right now)
		movie := c.FormValue("movie")

		// Use the Scraper
		moviesIds, err := Scraper(movie) // movie, maxusers, threads

		if err != nil {
			data["Error"] = "Movie not found"
			return c.Render(200, "error.html", data)
		}

		movieData = MovieData{movie, 0, moviesIds}

		// Use the LookUp to get the movie info
		// this also overwrites the movie data
		movies, _ = LookUpMovies(&movieData)

		// very rarely will be less than 10 movies
		if (movieData).IsFull() {
			data["HasMore"] = false
		}

		return c.Render(200, "recommendations.html", data)
	})

	e.POST("/loadMore", func(c echo.Context) error {
		// ensure we don't load movies if we don't have a movie loaded
		if len(movieData.ids) == 0 {
			data["Error"] = "No movies"
			return c.Render(200, "error.html", data)
		}

		TMDBMovieInfo, _ := LookUpMovies(&movieData)

		movies = append(movies, TMDBMovieInfo...) // concatenate the slices

		// will not show the load more button if there are no more movies
		if (movieData).IsFull() {
			data["HasMore"] = false
		}

		return c.Render(200, "recommendations.html", data)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
