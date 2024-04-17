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

func KeyedLookUpMovie(id int, key string) (Movie, error) {
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

	// if it unmarshals correctly, return the movie, otherwise return an error
	if movie.Id == 0 || movie.Title == "" || movie.Poster_path == "" || movie.Release_date == "" {
		return movie, fmt.Errorf("Movie ID not found on TMDB")
	}

	// get the year from the release date
	movie.Release_date = movie.Release_date[:4]

	return movie, nil
}

func LookUpMovie(movie string) (Movie, error) {
	// load the .env file
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("TMDB_API_KEY")

	// get the TMDB id
	id := convertToTMDBIds([]string{movie})[0]

	// look up the movie on TMDB
	movieData, err := KeyedLookUpMovie(id, key)

	return movieData, err
}

func LookUpMovies(movieData *MovieData) ([]Movie, error) {
	// load the .env file
	err := godotenv.Load()

	// do this if you're running the server locally
	// err := godotenv.Load("../.env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	key := os.Getenv("TMDB_API_KEY")

	movies := []Movie{}

	// convert the slugs to TMDB ids
	ids := convertToTMDBIds(movieData.slugs[movieData.pointer:movieData.Increment()])

	// look up each movie on TMDB
	for _, id := range ids {
		movie, err := KeyedLookUpMovie(id, key)

		// skip movies if they're not found
		if err == nil {
			movies = append(movies, movie)
		}
	}

	return movies, nil
}
