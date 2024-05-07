<div align="center">
  <h3 align="center">letterboxd-rec</h3>

  <p align="center">
    Accurate movie recommendations based on favourites on <a href="https://letterboxd.com/" target="_blank">Letterboxd</a>
    <br />
  </p>
</div>

<!-- ABOUT THE PROJECT -->
# About The Project

letterboxd-rec, provides *accurate* film recommendations provided only a single movie. These are calculated based on how many other users on Letterboxd have it listed on their profile, as one of their **top 4** movies.

Built entirely in Go using a go http web server, with [templ][Templ-URL] to help manage HTML templates and **HTMX**, this project is built on essentially **zero** javascript. To help with gathering users' *public* top 4's, it uses the web scraping framework [Colly][Colly-URL], built in Go. Movie data is gathered using TMDB's API.

## Rationale

Imagine you wanted to find how many other people in the world had the same favourite movie as you. If you asked them all what some of their *other* favourite movies are, chances are you would also enjoy those movies. If multiple people said the same movie multiple times, then that movie is probably even more likely to be one you enjoy - This is essentially the idea why I built this project, and why it can be very accurate at making good film recommendations.

There are some slight limitations, however; to get the most accurate results, you need to sprawl through a very large amount of users, which can take quite a while especially considering its built using a web scraper, which is necessary as Letterboxd's API is not public. Also, the most frequent recommendations tend to be analogous, because users' top 4 are skewed towards the most popular movies. So, if you are seeing Interstellar, Fight Club and La La Land for every recommendation, it's because they're some of the most popular movies!


### Features
* Movie Recommendations
* Film searching
* Responsive UI
* Tunable Parameters (users to scrape, server threads)


### Built With

* [templ][Templ-URL]
* [Colly][Colly-URL]

</br>

[![Go](https://img.shields.io/badge/Go-%2300ADD8.svg?&logo=go&logoColor=white)](#)
[![HTMX](https://img.shields.io/badge/%3C/%3E%20HTMX-3D72D7?logo=mysl&logoColor=white)](#)
[![TailwindCSS](https://img.shields.io/badge/Tailwind%20CSS-%2338B2AC.svg?logo=tailwind-css&logoColor=white)](#)

### Running the Webserver 

You should be able to simply clone the GitHub repository and run main.go, ensuring you have all the required dependencies which have been provided in the go.mod file.

Make sure you setup a .env file for the TMDB API keys, like this:


```
.env

TMDB_API_KEY="Your TMDB Access Token Auth"
API_KEY="Your TMDB API Key Auth"
```

<!-- ACKNOWLEDGMENTS -->
## Acknowledgments

* [Letterboxd][Letterboxd-URL]
* [TMDB][TMDB-URL]

## Contact

Jackson Lawrence - jplqqz@gmail.com

<!-- LINKS & IMAGES -->
[Templ-URL]: https://github.com/a-h/templ
[Colly-URL]: https://github.com/gocolly/colly
[Letterboxd-URL]: https://letterboxd.com/
[TMDB-URL]: https://www.themoviedb.org/
