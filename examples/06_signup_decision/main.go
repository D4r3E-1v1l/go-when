package main

import (
	"fmt"
	"strings"

	when "github.com/D4r3E-1v1l/go-when"
)

type SignupForm struct {
	Email string
	Age   int
	Plan  string
}

func main() {
	form := SignupForm{
		Email: "alice@example.com",
		Age:   16,
		Plan:  "pro",
	}

	decision := when.MatchAnyAs[string](form).
		When(isInvalidEmail).Then("invalid_email").
		When(needsParentApproval).Then("needs_parent_approval").
		When(isAllowedPlan).Then("allow").
		Else("manual_review")

	fmt.Println(decision)
}

func isInvalidEmail(form SignupForm) bool {
	return !strings.Contains(form.Email, "@")
}

func needsParentApproval(form SignupForm) bool {
	return form.Age < 18 && form.Plan == "pro"
}

func isAllowedPlan(form SignupForm) bool {
	return form.Plan == "free" || form.Plan == "pro"
}
