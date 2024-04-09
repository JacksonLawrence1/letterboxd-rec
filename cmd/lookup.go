package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Movie struct {
	Id           int
	Title        string
	Release_date string
	Poster_path  string
}

func LookUpMovie(id int, key string) (Movie, error) {
	movie := Movie{}

	// endpoint
	url := "https://api.themoviedb.org/3/movie/" + fmt.Sprint(id) + "?language=en-US"
	req, _ := http.NewRequest("GET", url, nil)

	// set the headers
	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return movie, err
	}

	// filter out the data we need from the json response
	defer res.Body.Close()

	// decode the json
	body, _ := io.ReadAll(res.Body)

	// only need the title, id and poster path (for now)
	json.Unmarshal(body, &movie)

	return movie, nil
}

func LookUpMovies(movieIds []int) ([]Movie, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("TMDB_API_KEY")

	movies := []Movie{}

	// look up each movie on TMDB
	for _, id := range movieIds {
		movie, err := LookUpMovie(id, key)

		if err != nil {
			return movies, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
