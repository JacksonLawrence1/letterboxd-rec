package main

import (
	"fmt"
	"slices"
	"strconv"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

type Movie struct {
	Title string
	Slug  string
	Id    string
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

func Scraper(movie string, maxUsersStr string, threadsStr string) ([]string, time.Duration) {
	start := time.Now()

	// convert strings to ints
	maxUsers, userErr := strconv.Atoi(maxUsersStr)
	threads, threadsErr := strconv.Atoi(threadsStr)

	if movie == "" || userErr != nil || threadsErr != nil {
		fmt.Println("Invalid inputs: ", movie, maxUsers, threads)
		return nil, time.Since(start)
	}

	// this users the movie's film slug, make sure you look up the correct one
	users := scrapeUsers(movie, maxUsers)

	// depending on how many users, this could take a while
	movies := scrapeFavourites(users, maxUsers, threads)

	keys := make([]string, 0, len(movies))
	for k := range movies {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, func(i string, j string) int {
		return movies[j] - movies[i]
	})

	return keys, time.Since(start)
}
