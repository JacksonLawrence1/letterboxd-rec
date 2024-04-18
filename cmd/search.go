package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/joho/godotenv"
)

type MovieSearchData struct {
	Results []Movie
}

func findSlugs(movies *map[int]Movie) []Movie {
	c := colly.NewCollector()

	q, _ := queue.New(
		Threads,
		&queue.InMemoryQueueStorage{MaxSize: len(*movies)},
	)

	c.OnHTML("body", func(e *colly.HTMLElement) {
		// split the url into parts separated by "/"
		split := strings.Split(e.Request.URL.String(), "/")
		size := len(split)

		// if the third to last part is "film", then we know letterboxd has the movie
		if size > 2 && split[size-3] == "film" {
			// on the page, get the tmdb id
			id, err := strconv.Atoi(e.Attr("data-tmdb-id"))

			// set the slug to movie
			if err == nil {
				movie := (*movies)[id]
				movie.Slug = split[size-2] // 2nd to last part is the slug
				if len(movie.Release_date) > 4 {
					movie.Release_date = movie.Release_date[:4] // only get the year
				}

				(*movies)[id] = movie
			}
		}
	})

	for _, movie := range *movies {
		url := "https://letterboxd.com/tmdb/" + fmt.Sprint(movie.Id) + "/"

		q.AddURL(url)
	}

	q.Run(c)

	// convert the map to a slice
	movieSlice := make([]Movie, 0, len(*movies))

	for _, movie := range *movies {
		movieSlice = append(movieSlice, movie)
	}

	return movieSlice
}

func search(term string) []Movie {
	// load the .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("TMDB_API_KEY")

	url := "https://api.themoviedb.org/3/search/movie?query=" + strings.ReplaceAll(term, " ", "%20") + "&include_adult=false&language=en-US&page=1"
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	body, _ := io.ReadAll(res.Body)

	searchResults := MovieSearchData{}
	json.Unmarshal(body, &searchResults)

	movieMap := map[int]Movie{}

	for _, movie := range searchResults.Results {
		if movie.Popularity > 10 {
			movieMap[movie.Id] = movie
		}
	}

	// scrape letterboxds to find its film slug

	searchResults.Results = findSlugs(&movieMap)

	// sort the movies by popularity
	//sort.Slice(searchResults.Results, func(i, j int) bool {
	//return searchResults.Results[i].Popularity > searchResults.Results[j].Popularity
	//})

	return searchResults.Results
}
