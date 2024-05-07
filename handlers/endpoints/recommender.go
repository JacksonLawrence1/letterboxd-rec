package endpoints

import (
	"letterboxd-rec/services"
	"letterboxd-rec/templates/pages"
	"letterboxd-rec/templates/partials"
	"letterboxd-rec/utils"
	"strconv"

	"net/http"

	"github.com/gorilla/schema"
)

func RecommendHandler(mux *http.ServeMux) {
	movies := []utils.Movie{}
	isFull := false

	// The user requests a recommendation from a movie
	mux.HandleFunc("POST /recommend", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()

		decoder := schema.NewDecoder()

		var movie utils.Movie
		err := decoder.Decode(&movie, r.PostForm)

		if err != nil {
			pages.ErrorPage(500, "Error while searching, please try again.").Render(r.Context(), w)
			return
		}

		// Clear current recommendations
		movies = []utils.Movie{}
		isFull = false

		utils.Progress = utils.ProgressData{Percent: 0, Message: "gathering users"}

		// asyncronously start recommendation process
		go func() {
			users, threads := utils.MaxUsers, utils.Threads

			// Get the recommendations for the selected movie
			movies, isFull, err = services.Recommend(movie, users, threads)

			if err != nil {
				pages.ErrorPage(500, err.Error()).Render(r.Context(), w)
			}

			utils.Progress.Percent = 100
		}()

		loadScreen := partials.Loading(utils.Progress, movie.Title)
		loadScreen.Render(r.Context(), w)
	})

	// Request the progress of the recommendation
	mux.HandleFunc("GET /recommend/progress", func(w http.ResponseWriter, r *http.Request) {
		if utils.Progress.Percent != 100 {
			progress := partials.ProgressBar(utils.Progress)
			progress.Render(r.Context(), w)
		}

		w.Header().Set("HX-Trigger", "done")
	})

	// After recommendation process is done, show the recommendations
	mux.HandleFunc("GET /recommendations", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		movieTitle := r.FormValue("Title")

		// if we dont have 10 movies in our first batch, we need to make sure load more button is not shown
		recommendationPanel := partials.RecommendationPanel(movieTitle, movies, min(utils.ItemsToShow+1, len(movies)+1), isFull)
		recommendationPanel.Render(r.Context(), w)
	})

	// Load more recommendations
	mux.HandleFunc("POST /loadMore", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		pointer, _ := strconv.Atoi(r.FormValue("load-more"))

		TMDBMovieInfo, newPointer, isFull := services.GetMoreRecommendations(pointer, utils.Threads)

		// Add the new recommendations
		if len(TMDBMovieInfo) > 0 {
			updatedRecommendations := partials.Recommendations(TMDBMovieInfo, newPointer, isFull)
			updatedRecommendations.Render(r.Context(), w)
		}
	})
}
