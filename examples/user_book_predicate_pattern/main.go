package main

import (
	"fmt"
	"strings"

	when "github.com/D4r3E-1v1l/go-when"
)

type User struct {
	Name string
	Age  int
	Role string
}

type Book struct {
	Title     string
	Available bool
}

type BorrowRequest struct {
	User User
	Book Book
}

type HasRole struct {
	Role string
}

func (p HasRole) Match(req BorrowRequest) bool {
	return req.User.Role == p.Role
}

func main() {
	req := BorrowRequest{
		User: User{Name: "Alice", Age: 16, Role: "student"},
		Book: Book{Title: "Learning Go", Available: true},
	}

	decision := when.MatchAnyAs[string](req).
		When(func(req BorrowRequest) bool {
			return !req.Book.Available
		}).Then("book_unavailable").
		When(func(req BorrowRequest) bool {
			return req.User.Age < 12 && strings.Contains(req.Book.Title, "Advanced")
		}).Then("needs_parent_approval").
		Pattern(HasRole{Role: "student"}).Then("student_can_borrow").
		Else("manual_review")

	fmt.Println(decision)
}
