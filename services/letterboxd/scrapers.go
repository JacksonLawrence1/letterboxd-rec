package letterboxd

import (
	"fmt"
	"strconv"
	"strings"

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
	progress := 0

	c.OnHTML("td.table-person", func(e *colly.HTMLElement) {
		if len(users) >= utils.MaxUsers {
			return
		}

		users = append(users, e.ChildAttr("a", "href"))

		// Update progress
		progress++
		utils.Progress.Percent = progress * 20 / utils.MaxUsers
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
	progress := 0

	c.OnHTML("section#favourites", func(e *colly.HTMLElement) {
		// add the 4 favourites to the map
		e.ForEach("div.poster", func(_ int, e *colly.HTMLElement) {
			movies[e.Attr("data-film-slug")]++
		})

		// Update progress
		progress++
		utils.Progress.Percent = 20 + (progress * 70 / len(users))
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

	TMDBIds := make([]int, len(movieSlugs))
	orderMap := make(map[string]int)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// get the TMDB id from the data attribute
		id, err := strconv.Atoi(e.Attr("data-tmdb-id"))
		split := strings.Split(e.Request.URL.String(), "/")

		if err == nil && id != 0 {
			slug := split[len(split)-2]

			// this maintains the correct order of the movies
			TMDBIds[orderMap[slug]] = id
		}
	})

	// scrape the TMDB id from the movie's letterboxd page
	for i, slug := range movieSlugs {
		orderMap[slug] = i
		q.AddURL(fmt.Sprintf("https://letterboxd.com/film/%s/", slug))
	}

	q.Run(c)

	return TMDBIds
}
