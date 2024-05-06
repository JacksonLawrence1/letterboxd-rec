package endpoints

import (
	"letterboxd-rec/services"
	"letterboxd-rec/templates/pages"
	"letterboxd-rec/templates/partials"
	"letterboxd-rec/utils"

	"net/http"

	"github.com/gorilla/schema"
)

func RecommendHandler(mux *http.ServeMux) {
	movies := []utils.Movie{}

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

		utils.Progress = utils.ProgressData{Percent: 0, Message: "gathering users"}

		go func() {
			// Get the recommendations for the selected movie
			movies, err = services.Recommend(movie)

			if err != nil {
				pages.ErrorPage(500, err.Error()).Render(r.Context(), w)
			}

			utils.Progress.Percent = 100
		}()

		loadScreen := partials.Loading(utils.Progress, movie.Title)
		loadScreen.Render(r.Context(), w)
	})

	mux.HandleFunc("GET /recommend/progress", func(w http.ResponseWriter, r *http.Request) {
		if utils.Progress.Percent != 100 {
			progress := partials.ProgressBar(utils.Progress)
			progress.Render(r.Context(), w)
		}

		w.Header().Set("HX-Trigger", "done")
	})

	mux.HandleFunc("GET /recommendations", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		movieTitle := r.FormValue("Title")

		recommendationPanel := partials.RecommendationPanel(movieTitle, movies)
		recommendationPanel.Render(r.Context(), w)
	})

	mux.HandleFunc("POST /isFull", func(w http.ResponseWriter, r *http.Request) {
		partials.LoadMore().Render(r.Context(), w)
	})

	mux.HandleFunc("POST /loadMore", func(w http.ResponseWriter, r *http.Request) {
		TMDBMovieInfo, isFull := services.GetMoreRecommendations()

		if isFull {
			w.Header().Set("HX-Trigger", "moreResults")
		}

		// Add the new recommendations
		if len(TMDBMovieInfo) > 0 {
			updatedRecommendations := partials.Recommendations(TMDBMovieInfo)
			updatedRecommendations.Render(r.Context(), w)
		}
	})
}
