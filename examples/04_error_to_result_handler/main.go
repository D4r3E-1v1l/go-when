package main

import (
	"errors"
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrInvalidInput  = errors.New("invalid input")
	ErrNotAuthorized = errors.New("not authorized")
)

type APIResult struct {
	Status  int
	Message string
}

type ErrorHandler func(error) APIResult

func main() {
	err := ErrInvalidInput

	handler := when.Err[ErrorHandler](err).
		Nil().Then(success).
		Is(ErrUserNotFound).Then(notFound).
		Is(ErrInvalidInput).Then(badRequest).
		Is(ErrNotAuthorized).Then(unauthorized).
		Else(internalError)

	result := handler(err)

	fmt.Printf("%d %s\n", result.Status, result.Message)
}

func success(error) APIResult {
	return APIResult{
		Status:  200,
		Message: "success",
	}
}

func notFound(error) APIResult {
	return APIResult{
		Status:  404,
		Message: "user not found",
	}
}

func badRequest(error) APIResult {
	return APIResult{
		Status:  400,
		Message: "invalid input",
	}
}

func unauthorized(error) APIResult {
	return APIResult{
		Status:  401,
		Message: "not authorized",
	}
}

func internalError(err error) APIResult {
	return APIResult{
		Status:  500,
		Message: "internal error: " + err.Error(),
	}
}
