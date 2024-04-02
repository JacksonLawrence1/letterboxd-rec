package main

import (
	"fmt"
	"slices"

	"github.com/gocolly/colly"
)

// gets the usernames of user's who have this movie in their top 4
func scrapeUsers(movie string, pages int, maxUsers int) []string {
	c := colly.NewCollector()

	var users []string

	c.OnHTML("td.table-person", func(e *colly.HTMLElement) {
		if maxUsers != -1 && len(users) >= maxUsers {
			fmt.Print(maxUsers, len(users))
			return
		}

		users = append(users, e.ChildAttr("a", "href"))
	})

	// letterboxd sorts users in alphabetical order - depending on you, this could be a good or bad thing
	for i := 1; i <= pages; i++ {
		if maxUsers != -1 && len(users) >= maxUsers {
			break
		}
		c.Visit(fmt.Sprintf("https://letterboxd.com/film/%s/fans/page/%d/", movie, i))
	}

	return users
}

// for each user, add their 4 favourites to a map
func scrapeFavourites(users []string) map[string]int {
	c := colly.NewCollector()

	// now we have the users, go onto their pages and scrape their 4 favourites (excluding the movie we're searching for)
	movies := make(map[string]int) // map movie's slug to number of users who like it

	c.OnHTML("section#favourites", func(e *colly.HTMLElement) {
		// add the 4 favourites to the map
		e.ForEach("div.poster", func(_ int, e *colly.HTMLElement) {
			movies[e.Attr("data-film-slug")]++
		})
	})

	for _, user := range users {
		c.Visit(fmt.Sprintf("https://letterboxd.com%s", user))
	}

	return movies
}

func Scraper() {
	// this users the movie's film slug, make sure you look up the correct one
	users := scrapeUsers("parasite-2019", 1, -1)

	// depending on how many users, this could take a while
	movies := scrapeFavourites(users)

	keys := make([]string, 0, len(movies))
	for k := range movies {
		keys = append(keys, k)
	}

	slices.SortFunc(keys, func(i string, j string) int {
		return movies[j] - movies[i]
	})

	i := 0

	// print out all movies that have more than 1 count
	for movies[keys[i]] > 1 {
		fmt.Printf("%s: %d\n", keys[i], movies[keys[i]])
		i++
	}
}
