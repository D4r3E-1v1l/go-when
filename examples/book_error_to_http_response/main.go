package main

import (
	"errors"
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

var (
	ErrBookNotFound      = errors.New("book not found")
	ErrBookAlreadyExists = errors.New("book already exists")
)

type HTTPResponse struct {
	Status int
	Body   string
}

type ValidationError struct {
	Field string
	Msg   string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Msg
}

func main() {
	err := fmt.Errorf("create book failed: %w", &ValidationError{
		Field: "title",
		Msg:   "must not be empty",
	})

	resp := mapErrorToHTTPResponse(err)
	fmt.Printf("%d %s\n", resp.Status, resp.Body)
}

func mapErrorToHTTPResponse(err error) HTTPResponse {
	var validationErr *ValidationError

	return when.Err[HTTPResponse](err).
		Nil().Then(HTTPResponse{Status: 200, Body: "ok"}).
		Is(ErrBookNotFound).Then(HTTPResponse{Status: 404, Body: "book not found"}).
		Is(ErrBookAlreadyExists).Then(HTTPResponse{Status: 409, Body: "book already exists"}).
		As(&validationErr).ThenDo(func(error) HTTPResponse {
		return HTTPResponse{
			Status: 400,
			Body:   "invalid " + validationErr.Field + ": " + validationErr.Msg,
		}
	}).
		Contains("timeout").Then(HTTPResponse{Status: 504, Body: "gateway timeout"}).
		Else(HTTPResponse{Status: 500, Body: "internal server error"})
}
