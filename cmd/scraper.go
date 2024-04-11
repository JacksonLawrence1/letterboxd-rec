package main

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

// takes the total users for the movie and randomises which pages to scrape
func randomise(totalUsers int, targetPages int) []int {
	// based on the total users, calculate the pages it has
	pages := totalUsers/25 + 1

	// select targetPages randomly between 1 and pages (no duplicates)
	var randomPages = make([]int, pages) // TODO

	return randomPages
}

// gets the usernames of user's who have this movie in their top 4
func scrapeUsers(movie string, maxUsers int, threads int) []string {
	c := colly.NewCollector()

	users := []string{}

	c.OnHTML("td.table-person", func(e *colly.HTMLElement) {
		if maxUsers != -1 && len(users) >= maxUsers {
			return
		}

		users = append(users, e.ChildAttr("a", "href"))
	})

	pages := maxUsers/25 + 1

	q, _ := queue.New(
		threads,
		&queue.InMemoryQueueStorage{MaxSize: pages},
	)

	// letterboxd sorts users in alphabetical order - depending on you, this could be a good or bad thing

	// future improvement: randomise the pages

	for i := 1; i <= pages; i++ {
		if maxUsers != -1 && len(users) >= maxUsers {
			break
		}
		q.AddURL(fmt.Sprintf("https://letterboxd.com/film/%s/fans/page/%d/", movie, i))
	}

	q.Run(c)

	return users
}

// for each user, add their 4 favourites to a map
func scrapeFavourites(users []string, maxUsers int, threads int) map[string]int {
	c := colly.NewCollector()

	if maxUsers == -1 {
		maxUsers = len(users)
	}

	// create a queue with a worker pool of 2 threads
	q, _ := queue.New(
		threads,
		&queue.InMemoryQueueStorage{MaxSize: maxUsers},
	)

	// now we have the users, go onto their pages and scrape their 4 favourites (excluding the movie we're searching for)
	movies := make(map[string]int) // map movie's slug to number of users who like it

	c.OnHTML("section#favourites", func(e *colly.HTMLElement) {
		// add the 4 favourites to the map
		e.ForEach("div.poster", func(_ int, e *colly.HTMLElement) {
			movies[e.Attr("data-film-slug")]++
		})
	})

	// ideally, we want to only store the ids, then fetch all info we need from TMDB

	for _, user := range users {
		q.AddURL(fmt.Sprintf("https://letterboxd.com%s", user))
	}

	q.Run(c)

	return movies
}

func convertToTMDBIds(movieSlugs []string, threads int) []int {
	c := colly.NewCollector()

	q, _ := queue.New(
		threads,
		&queue.InMemoryQueueStorage{MaxSize: len(movieSlugs)},
	)

	TMDBIds := []int{}

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// get the TMDB id from the data attribute
		id, err := strconv.Atoi(e.Attr("data-tmdb-id"))

		if err == nil && id != 0 {
			TMDBIds = append(TMDBIds, id)
		}
	})

	// scrape the TMDB id from the movie's letterboxd page
	for _, id := range movieSlugs {
		q.AddURL(fmt.Sprintf("https://letterboxd.com/film/%s/", id))
	}

	q.Run(c)

	return TMDBIds
}

func MovieExists(movie string) bool {
	c := colly.NewCollector()

	exists := false

	c.OnHTML("body", func(e *colly.HTMLElement) {
		exists = true
	})

	c.Visit(fmt.Sprintf("https://letterboxd.com/film/%s/", movie))
	return exists
}

func Scraper(movie string, maxUsers int, threads int) ([]int, error) {
	// ensure the movie exists on letterboxd
	if !MovieExists(movie) {
		return nil, fmt.Errorf("Movie not found on Letterboxd")
	}

	// this users the movie's film slug, make sure you look up the correct one
	users := scrapeUsers(movie, maxUsers, threads)

	// depending on how many users, this could take a while
	movies := scrapeFavourites(users, maxUsers, threads)

	keys := make([]string, 0, len(movies))
	for k := range movies {
		keys = append(keys, k)
	}

	// sort the movies by the number of users who like it in descending order
	slices.SortFunc(keys, func(i string, j string) int {
		return movies[j] - movies[i]
	})

	// because letterboxd doesn't store the TMDB id on the favourites page, scrape it from the movie details page
	ids := convertToTMDBIds(keys, threads)

	return ids, nil
}
