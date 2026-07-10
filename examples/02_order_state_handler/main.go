package main

import (
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

type OrderState int

const (
	OrderCreated OrderState = iota
	OrderPaid
	OrderPacked
	OrderShipped
	OrderCancelled
)

type Handler func()

func main() {
	state := OrderPaid

	handler := when.MatchAs[Handler](state).
		Case(OrderCreated).Then(requestPayment).
		Case(OrderPaid).Then(packOrder).
		Case(OrderPacked).Then(shipOrder).
		Case(OrderShipped).Then(sendTracking).
		Case(OrderCancelled).Then(cancelOrder).
		Exhaustive()

	handler()
}

func requestPayment() {
	fmt.Println("request payment")
}

func packOrder() {
	fmt.Println("pack order")
}

func shipOrder() {
	fmt.Println("ship order")
}

func sendTracking() {
	fmt.Println("send tracking information")
}

func cancelOrder() {
	fmt.Println("cancel order")
}
