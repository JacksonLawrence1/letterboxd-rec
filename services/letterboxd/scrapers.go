package letterboxd

import (
	"fmt"

	"letterboxd-rec/utils"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func newScraper(maxSize int, threads int) (*colly.Collector, *queue.Queue) {
	c := colly.NewCollector()

	q, _ := queue.New(
		threads,
		&queue.InMemoryQueueStorage{MaxSize: maxSize},
	)

	return c, q
}

// gets the usernames of user's who have this movie in their top 4
func ScrapeUsers(movie *utils.Movie, maxUsers int, threads int) []string {
	c, q := newScraper(maxUsers, threads)

	users := []string{}
	progress := 0

	c.OnHTML("td.table-person", func(e *colly.HTMLElement) {
		if len(users) >= maxUsers {
			return
		}

		users = append(users, e.ChildAttr("a", "href"))

		// Update progress
		progress++
		utils.Progress.Percent = progress * 20 / maxUsers
	})

	pages := utils.RandomisePages(movie.Fans, maxUsers)

	// gets random pages based on the total number of users
	for key := range pages {
		q.AddURL(fmt.Sprintf("https://letterboxd.com/film/%s/fans/page/%d/", movie.Slug, key))
	}

	q.Run(c)

	return users
}

// for each user, add their 4 favourites to a map
func ScrapeFavourites(users []string, threads int) map[string]int {
	size := len(users)
	c, q := newScraper(size, threads)

	// now we have the users, go onto their pages and scrape their 4 favourites (excluding the movie we're searching for)
	movies := make(map[string]int) // map movie's slug to number of users who like it
	progress := 0

	c.OnHTML("section#favourites", func(e *colly.HTMLElement) {
		// add the 4 favourites to the map
		e.ForEach("div.poster", func(_ int, e *colly.HTMLElement) {
			// check if e.Attr("data-film-slug") exists on the page
			if e.Attr("data-film-slug") != "" {
				movies[e.Attr("data-film-slug")]++
			}
		})

		// Update progress
		progress++
		utils.Progress.Percent = 20 + (progress * 70 / size)
	})

	// ideally, we want to only store the ids, then fetch all info we need from TMDB

	for _, user := range users {
		q.AddURL(fmt.Sprintf("https://letterboxd.com%s", user))
	}

	q.Run(c)

	return movies
}
