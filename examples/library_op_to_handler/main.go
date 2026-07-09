package main

import (
	"errors"
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

type LibraryOp int

const (
	BorrowBook LibraryOp = iota
	ReturnBook
	ReserveBook
	DeleteBook
)

var (
	ErrUnsupportedOperation = errors.New("unsupported library operation")
	ErrDeleteDisabled       = errors.New("delete book operation is disabled")
)

type Request struct {
	UserName string
	BookName string
}

type Response struct {
	Message string
}

type Handler func(Request) Response

func main() {
	op := BorrowBook
	req := Request{
		UserName: "Alice",
		BookName: "The Go Programming Language",
	}

	handler, err := when.MatchAs[Handler](op).
		WithErr().
		Case(BorrowBook).Then(borrowBook).
		Case(ReturnBook).Then(returnBook).
		Case(ReserveBook).Then(reserveBook).
		Case(DeleteBook).ThenErr(nil, ErrDeleteDisabled).
		ElseErr(nil, ErrUnsupportedOperation)

	if err != nil {
		fmt.Println("dispatch failed:", err)
		return
	}

	resp := handler(req)
	fmt.Println(resp.Message)
}

func borrowBook(req Request) Response {
	return Response{
		Message: req.UserName + " borrowed " + req.BookName,
	}
}

func returnBook(req Request) Response {
	return Response{
		Message: req.UserName + " returned " + req.BookName,
	}
}

func reserveBook(req Request) Response {
	return Response{
		Message: req.UserName + " reserved " + req.BookName,
	}
}
