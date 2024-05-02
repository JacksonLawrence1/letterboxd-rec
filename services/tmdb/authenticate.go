package tmdb

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var key string = ""

// helper function to make requests to TMDB
func request_TMDB(url string) (*http.Response, error) {
	// don't want to load the .env file every time
	if key == "" {
		// load the .env file
		err := godotenv.Load()

		// do this if you're running the server locally
		// err := godotenv.Load("../.env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}

		key = os.Getenv("TMDB_API_KEY")
	}

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+key)

	return http.DefaultClient.Do(req)
}
