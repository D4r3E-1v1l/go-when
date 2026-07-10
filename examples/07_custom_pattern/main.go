package main

import (
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

type Cart struct {
	MemberLevel string
	Total       int
	CouponCode  string
}

type VIPCart struct{}

func (VIPCart) Match(cart Cart) bool {
	return cart.MemberLevel == "vip"
}

type MinimumTotal struct {
	Amount int
}

func (p MinimumTotal) Match(cart Cart) bool {
	return cart.Total >= p.Amount
}

func main() {
	cart := Cart{
		MemberLevel: "vip",
		Total:       120,
		CouponCode:  "SUMMER",
	}

	discount := when.MatchAnyAs[int](cart).
		Pattern(VIPCart{}).Then(20).
		Pattern(MinimumTotal{Amount: 100}).Then(10).
		Else(0)

	fmt.Println(discount)
}
