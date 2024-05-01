package services

import (
	"fmt"
	"math/rand"
	"slices"
	"strconv"
	"strings"

	"letterboxd-rec/utils"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

// takes the total users for the movie and randomises which pages to scrape
func randomise(fans int) map[int]bool {
	maxFans := min(256, fans/25) // letterboxd seems to have a maximum of 256 pages

	pages := max(1, utils.MaxUsers/25)

	// select targetPages randomly between 1 and pages (no duplicates)
	randomPages := make(map[int]bool)

	for i := 0; i < pages; i++ {
		num := rand.Intn(maxFans) + 1
		if !randomPages[num] {
			randomPages[num] = true
		}
	}

	return randomPages
}

// gets the usernames of user's who have this movie in their top 4
func scrapeUsers(movie string, fans int) []string {
	c := colly.NewCollector()

	users := []string{}

	c.OnHTML("td.table-person", func(e *colly.HTMLElement) {
		if len(users) >= utils.MaxUsers {
			return
		}

		users = append(users, e.ChildAttr("a", "href"))
	})

	pages := randomise(fans)
	fmt.Println(pages)

	q, _ := queue.New(
		utils.Threads,
		&queue.InMemoryQueueStorage{MaxSize: len(pages) + 1},
	)

	// gets random pages based on the total number of users
	for key := range pages {
		q.AddURL(fmt.Sprintf("https://letterboxd.com/film/%s/fans/page/%d/", movie, key))
	}

	q.Run(c)

	return users
}

// for each user, add their 4 favourites to a map
func scrapeFavourites(users []string) map[string]int {
	c := colly.NewCollector()

	// create a queue with a worker pool of 2 threads
	q, _ := queue.New(
		utils.Threads,
		&queue.InMemoryQueueStorage{MaxSize: utils.MaxUsers},
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

func convertToTMDBIds(movieSlugs []string) []int {
	c := colly.NewCollector()

	q, _ := queue.New(
		utils.Threads,
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

func MovieExists(movie string) (bool, int) {
	c := colly.NewCollector()

	exists := false

	c.OnHTML("body", func(e *colly.HTMLElement) {
		exists = true
	})

	fans := 0

	c.OnHTML("li.js-route-fans", func(e *colly.HTMLElement) {
		fanString := e.ChildAttr("a", "title")
		fanString = strings.ReplaceAll(fanString[:len(fanString)-6], ",", "") // remove the "fans" part
		fans, _ = strconv.Atoi(fanString)
	})

	c.Visit(fmt.Sprintf("https://letterboxd.com/film/%s/fans/", movie))
	return exists, fans
}

func Scraper(movie string) ([]string, error) {
	exists, fans := MovieExists(movie)

	// ensure the movie exists on letterboxd
	if !exists {
		return nil, fmt.Errorf("movie not found on letterboxd")
	}

	// this users the movie's film slug, make sure you look up the correct one
	users := scrapeUsers(movie, fans)

	// depending on how many users, this could take a while
	movies := scrapeFavourites(users)

	keys := make([]string, 0, len(movies))
	for k := range movies {
		// don't include the movie we're searching for
		if k != movie {
			keys = append(keys, k)
		}
	}

	// sort the movies by the number of users who like it in descending order
	slices.SortFunc(keys, func(i string, j string) int {
		return movies[j] - movies[i]
	})

	return keys, nil
}
