package utils

import (
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

var Threads int = 8
var MaxUsers int = 10
var ItemsToShow int = 10

type Movie struct {
	Id           int
	Title        string
	Slug         string
	Release_date string
	Poster_path  string
	Overview     string
	Popularity   float64
	Fans         int
}

type MovieData struct {
	Movie   Movie
	Pointer int
	Slugs   []string
}

func (s *MovieData) Increment() int {
	s.Pointer = min(s.Pointer+ItemsToShow, len(s.Slugs)) // how many movies at a time
	return s.Pointer
}

func (s *MovieData) IsFull() bool {
	return s.Pointer == len(s.Slugs)
}

func SortByFrequency(slug string, movies map[string]int) []string {
	keys := make([]string, 0, len(movies))

	for key := range movies {
		if key != slug {
			keys = append(keys, key)
		}
	}

	// sort the movies by the number of users who like it in descending order
	slices.SortFunc(keys, func(i string, j string) int {
		return movies[j] - movies[i]
	})

	return keys
}

func FilterByPopularity(movies []Movie, popularity float64) []Movie {
	filteredMovies := make([]Movie, 0)

	for _, movie := range movies {
		if movie.Popularity >= popularity {
			filteredMovies = append(filteredMovies, movie)
		}
	}

	return filteredMovies
}

// takes the total users for the movie and randomises which pages to scrape
func RandomisePages(fans int) map[int]bool {
	maxFans := min(256, fans/25) // letterboxd seems to have a maximum of 256 pages

	pages := max(1, MaxUsers/25)

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

func ConvertFans(fanString string) int {
	fanString = strings.ReplaceAll(fanString, ",", "")

	// first instance of '\u00a0' i.e. the whitespace
	i := strings.Index(fanString, "\u00a0")

	if i == -1 {
		return -1
	}

	res, err := strconv.Atoi(fanString[:i])

	if err != nil {
		return -1
	}

	return res
}
