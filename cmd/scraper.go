package main

import (
	"fmt"
	"slices"
	"strconv"
	"time"

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
func scrapeUsers(movie string, maxUsers int) []string {
	c := colly.NewCollector()

	var users []string

	c.OnHTML("td.table-person", func(e *colly.HTMLElement) {
		if maxUsers != -1 && len(users) >= maxUsers {
			return
		}

		users = append(users, e.ChildAttr("a", "href"))
	})

	pages := maxUsers/25 + 1

	// letterboxd sorts users in alphabetical order - depending on you, this could be a good or bad thing

	// future improvement: randomise the pages

	for i := 1; i <= pages; i++ {
		if maxUsers != -1 && len(users) >= maxUsers {
			break
		}
		c.Visit(fmt.Sprintf("https://letterboxd.com/film/%s/fans/page/%d/", movie, i))
	}

	return users
}

// for each user, add their 4 favourites to a map
func scrapeFavourites(users []string, maxUsers int, threads int) map[int]int {
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
	movies := make(map[int]int) // map movie's slug to number of users who like it

	c.OnHTML("section#favourites", func(e *colly.HTMLElement) {
		// add the 4 favourites to the map
		e.ForEach("div.poster", func(_ int, e *colly.HTMLElement) {
			id, err := strconv.Atoi(e.Attr("data-film-id"))
			if err == nil {
				movies[id]++
			}
		})
	})

	// ideally, we want to only store the ids, then fetch all info we need from TMDB

	for _, user := range users {
		q.AddURL(fmt.Sprintf("https://letterboxd.com%s", user))
	}

	q.Run(c)

	return movies
}

func Scraper(movie string, maxUsers int, threads int) ([]int, time.Duration) {
	start := time.Now()

	// this users the movie's film slug, make sure you look up the correct one
	users := scrapeUsers(movie, maxUsers)

	// depending on how many users, this could take a while
	movies := scrapeFavourites(users, maxUsers, threads)

	keys := make([]int, 0, len(movies))
	for k := range movies {
		keys = append(keys, k)
	}

	// sort the movies by the number of users who like it in descending order
	slices.SortFunc(keys, func(i int, j int) int {
		return movies[j] - movies[i]
	})

	return keys, time.Since(start)
}
