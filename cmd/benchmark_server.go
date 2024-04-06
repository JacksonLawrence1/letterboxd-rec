package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type thread struct {
	Count int
}

type Time struct {
	Time time.Duration
}

func StartBenchmarkServer() {
	e := echo.New()

	// e.Use(middleware.Logger())

	data := map[string]interface{}{}

	e.Renderer = newTemplate()

	// render index.html
	e.GET("/", func(c echo.Context) error {
		data["Movies"] = []Movie{
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
		return c.Render(200, "benchmark.html", data)
	})

	e.POST("/search", func(c echo.Context) error {
		// TODO: introduce some rate limiting

		chosen := c.FormValue("movie")
		maxUsersStr := c.FormValue("users")
		threadsStr := c.FormValue("threads")

		// convert strings to ints
		maxUsers, userErr := strconv.Atoi(maxUsersStr)
		threads, threadsErr := strconv.Atoi(threadsStr)

		if chosen == "" || userErr != nil || threadsErr != nil {
			fmt.Println("Invalid inputs: ", chosen, maxUsers, threads)
			return nil
		}

		fmt.Println("Starting search for", chosen, "with", maxUsers, "users and", threads, "threads")

		movies, time := Scraper(chosen, maxUsers, threads)

		fmt.Println(time)

		// map the first 5 movies to movie structs
		data["ReturnedMovies"] = []Movie{}

		// right now it shows the movie you chose, but we should remove it
		for i := 0; i < min(5, len(movies)); i++ {
			data["ReturnedMovies"] = append(data["ReturnedMovies"].([]Movie), Movie{Title: movies[i], Slug: movies[i]})
		}

		data["Time"] = Time{Time: time}

		// ideally here we render what movies we found
		return c.Render(200, "benchmark_results.html", data)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
