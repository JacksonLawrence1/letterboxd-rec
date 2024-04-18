package main

import (
	"github.com/labstack/echo/v4"
)

func movies(e *echo.Echo) {
	data := map[string]interface{}{}

	// essentially the "backend" data of the movies
	moviesData := MovieData{Movie{}, 0, []string{}}
	data["Movie"] = &moviesData.Movie

	// the frontend data of the movies, e.g. full title, release date and poster
	movies := []Movie{}
	data["Movies"] = &movies

	// whether there are more movies to load
	data["HasMore"] = true

	// what to run when the user searches for a movie
	e.POST("/search", func(c echo.Context) error {
		// use TMBD API to search for the movies
		movie := c.FormValue("movie")

		data["Term"] = movie
		data["Results"] = search(movie)

		return c.Render(200, "search_results", data)
	})

	// when we have the movie with exact data, search for recommendations
	e.POST("/recommend", func(c echo.Context) error {
		movie := c.FormValue("movie")

		data["HasMore"] = true

		// get movie data
		movieData, err := LookUpMovie(movie)

		if err != nil {
			return c.String(404, "Movie not found on TMDB")
		}

		// Use the Scraper to get the movie slugs
		movieSlugs, err := Scraper(movie) // movie, maxusers, threads

		if err != nil {
			return c.String(404, "Movie not found on Letterboxd")
		}

		moviesData = MovieData{movieData, 0, movieSlugs}

		// Use the LookUp to get the movie info
		// this also overwrites the movie data
		movies, _ = LookUpMovies(&moviesData)

		if len(movies) == 0 {
			return c.String(404, "Error when looking up movies")
		}

		// very rarely will be less than 10 movies
		if (moviesData).IsFull() {
			data["HasMore"] = false
		}

		return c.Render(200, "recommendations.html", data)
	})

	e.POST("/loadMore", func(c echo.Context) error {
		// ensure we don't load movies if we don't have a movie loaded
		if len(moviesData.slugs) == 0 {
			return c.String(404, "No movies")
		}

		TMDBMovieInfo, _ := LookUpMovies(&moviesData)

		movies = append(movies, TMDBMovieInfo...) // concatenate the slices

		// will not show the load more button if there are no more movies
		if (moviesData).IsFull() {
			data["HasMore"] = false
		}

		return c.Render(200, "recommendations.html", data)
	})
}
