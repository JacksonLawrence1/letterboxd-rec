package letterboxd

import (
	"fmt"
	"letterboxd-rec/utils"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

// Find only the movies that have more than maxUsers fans on letterboxd
func FilterLetterboxdMoviesByFans(movies []utils.Movie, maxUsers int) []utils.Movie {
	// very unlikely a movie with low tmdb popularity has enough fans on letterboxd
	movies = utils.FilterByPopularity(movies, 10)

	// filter movies that actually exist on letterboxd
	movies = FilterLetterboxdMovies(movies)

	filteredMovies := make([]utils.Movie, 0)
	slugMap := make(map[string]*utils.Movie)

	// avoid creating another scraper if there are no movies
	if len(movies) == 0 {
		return filteredMovies
	}

	c, q := newScraper(utils.Threads, len(movies))

	// get the number of fans for each movie
	c.OnHTML("li.js-route-fans", func(e *colly.HTMLElement) {
		fanString := e.ChildAttr("a", "title")

		// format the string and convert it to an integer
		fans := utils.ConvertFans(fanString)

		if fans >= maxUsers {
			slug := strings.Split(e.Request.URL.Path, "/")[2]

			if slugMap[slug] == nil {
				return
			}

			slugMap[slug].Fans = fans
			filteredMovies = append(filteredMovies, *slugMap[slug])
		}
	})

	for _, movie := range movies {
		slugMap[movie.Slug] = &movie

		q.AddURL("https://letterboxd.com/film/" + movie.Slug + "/fans/")
	}

	q.Run(c)

	return filteredMovies
}

// Given a list of movies, find only the ones that are available on letterboxd and append their slugs
func FilterLetterboxdMovies(movies []utils.Movie) []utils.Movie {
	c, q := newScraper(utils.Threads, len(movies))

	// map tmdb id to index in the movies slice
	tmdbMap := make(map[int]int)
	foundMovies := make([]utils.Movie, 0)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		split := strings.Split(e.Request.URL.String(), "/")
		size := len(split)

		// if the third to last part is "film", then it exists on letterboxd
		if size > 2 && split[size-3] == "film" {
			// on the page, get the tmdb id
			id, err := strconv.Atoi(e.Attr("data-tmdb-id"))

			if err != nil {
				return
			}

			// get the movie from the slice
			movie := &movies[tmdbMap[id]]
			movie.Slug = split[size-2] // 2nd to last part is the slug

			foundMovies = append(foundMovies, *movie)
		}
	})

	for i, movie := range movies {
		tmdbMap[movie.Id] = i

		url := "https://letterboxd.com/tmdb/" + fmt.Sprint(movie.Id) + "/"
		q.AddURL(url)
	}

	q.Run(c)
	return foundMovies
}
