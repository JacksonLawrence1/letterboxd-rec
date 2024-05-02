package letterboxd

import (
	"fmt"
	"strconv"

	"letterboxd-rec/utils"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
)

func newScraper(threads int, maxSize int) (*colly.Collector, *queue.Queue) {
	c := colly.NewCollector()

	q, _ := queue.New(
		threads,
		&queue.InMemoryQueueStorage{MaxSize: maxSize},
	)

	return c, q
}

// gets the usernames of user's who have this movie in their top 4
func ScrapeUsers(movie *utils.Movie) []string {
	c, q := newScraper(utils.Threads, utils.MaxUsers)

	users := []string{}

	c.OnHTML("td.table-person", func(e *colly.HTMLElement) {
		if len(users) >= utils.MaxUsers {
			return
		}

		users = append(users, e.ChildAttr("a", "href"))
	})

	pages := utils.RandomisePages(movie.Fans)

	// gets random pages based on the total number of users
	for key := range pages {
		q.AddURL(fmt.Sprintf("https://letterboxd.com/film/%s/fans/page/%d/", movie.Slug, key))
	}

	q.Run(c)

	return users
}

// for each user, add their 4 favourites to a map
func ScrapeFavourites(users []string) map[string]int {
	c, q := newScraper(utils.Threads, len(users))

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

// converts the movie slugs to TMDB ids so we can get the movie info from TMDB
func ConvertMovieSlugs(movieSlugs []string) []int {
	c, q := newScraper(utils.Threads, len(movieSlugs))

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
