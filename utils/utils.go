package utils

var Threads int = 8
var MaxUsers int = 5

type Movie struct {
	Id           int
	Title        string
	Slug         string
	Release_date string
	Poster_path  string
	Overview     string
	Popularity   float64
}

type MovieData struct {
	Movie   Movie
	Pointer int
	Slugs   []string
}

var ItemsToShow int = 10

func (s *MovieData) Increment() int {
	s.Pointer = min(s.Pointer+ItemsToShow, len(s.Slugs)) // how many movies at a time
	return s.Pointer
}

func (s *MovieData) IsFull() bool {
	return s.Pointer == len(s.Slugs)
}
