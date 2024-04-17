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
	moviesData := MovieData{Movie{}, 0, []string{}}
	data["Movie"] = &moviesData.Movie

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

		// use TMBD API to search for the movies
		movie := c.FormValue("movie")

		// get movie data
		movieData, err := LookUpMovie(movie)

		if err != nil {
			data["Error"] = "Movie not found on TMDB"
			return c.Render(200, "error.html", data)
		}

		// Use the Scraper to get the movie slugs
		movieSlugs, err := Scraper(movie) // movie, maxusers, threads

		if err != nil {
			data["Error"] = "Movie not found on Letterboxd"
			return c.Render(200, "error.html", data)
		}

		moviesData = MovieData{movieData, 0, movieSlugs}

		// Use the LookUp to get the movie info
		// this also overwrites the movie data
		movies, _ = LookUpMovies(&moviesData)

		// very rarely will be less than 10 movies
		if (moviesData).IsFull() {
			data["HasMore"] = false
		}

		return c.Render(200, "recommendations.html", data)
	})

	e.POST("/loadMore", func(c echo.Context) error {
		// ensure we don't load movies if we don't have a movie loaded
		if len(moviesData.slugs) == 0 {
			data["Error"] = "No movies"
			return c.Render(200, "error.html", data)
		}

		TMDBMovieInfo, _ := LookUpMovies(&moviesData)

		movies = append(movies, TMDBMovieInfo...) // concatenate the slices

		// will not show the load more button if there are no more movies
		if (moviesData).IsFull() {
			data["HasMore"] = false
		}

		return c.Render(200, "recommendations.html", data)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
