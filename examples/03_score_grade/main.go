package main

import (
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

func main() {
	score := 86

	grade := when.MatchAs[string](score).
		Case(100).Then("A+").
		Range(when.Range(90, 100)).Then("A").
		Range(when.Range(80, 90)).Then("B").
		Range(when.Range(70, 80)).Then("C").
		Range(when.Range(60, 70)).Then("D").
		Range(when.Range(0, 60)).Then("F").
		Else("Invalid score")

	fmt.Println(grade)
}
