package main

var ItemsToShow int = 10
var Threads int = 8
var MaxUsers int = 100

type Movie struct {
	Id           int
	Title        string
	Release_date string
	Poster_path  string
}

type MovieData struct {
	movie   string
	pointer int
	ids     []int
}

func (s *MovieData) Increment() int {
	s.pointer = min(s.pointer+ItemsToShow, len(s.ids)) // how many movies at a time
	return s.pointer
}

func (s *MovieData) IsFull() bool {
	return s.pointer == len(s.ids)
}
