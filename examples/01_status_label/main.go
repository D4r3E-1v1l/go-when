package main

import (
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

func main() {
	status := "paid"

	label := when.MatchAs[string](status).
		Case("pending").Then("Waiting for payment").
		Case("paid").Then("Payment received").
		Case("shipped").Then("On the way").
		Case("cancelled").Then("Order cancelled").
		Else("Unknown status")

	fmt.Println(label)
}
