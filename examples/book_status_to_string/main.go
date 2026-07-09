package main

import (
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

type BookStatus int

const (
	BookAvailable BookStatus = iota
	BookBorrowed
	BookReserved
	BookLost
)

func main() {
	status := BookReserved

	text := when.MatchAs[string](status).
		Case(BookAvailable).Then("available").
		Case(BookBorrowed).Then("borrowed").
		Case(BookReserved).Then("reserved").
		Case(BookLost).Then("lost").
		Else("unknown")

	fmt.Println(text)
}
