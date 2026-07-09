package main

import (
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

func main() {
	readingScore := 86

	grade := when.MatchAs[string](readingScore).
		Range(when.Range(0, 60)).Then("F").
		Range(when.Range(60, 70)).Then("D").
		Range(when.Range(70, 80)).Then("C").
		Range(when.Range(80, 90)).Then("B").
		Range(when.Closed(90, 100)).Then("A").
		Else("invalid")

	fmt.Println(grade)
}
