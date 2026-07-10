package main

import (
	"errors"
	"fmt"

	when "github.com/D4r3E-1v1l/go-when"
)

type Command string

const (
	AddTodo    Command = "add"
	DoneTodo   Command = "done"
	DeleteTodo Command = "delete"
	ExportTodo Command = "export"
)

var (
	ErrUnsupportedCommand = errors.New("unsupported command")
	ErrExportDisabled     = errors.New("export is disabled")
)

type Handler func()

func main() {
	command := ExportTodo

	handler, err := when.MatchAs[Handler](command).
		WithErr().
		Case(AddTodo).Then(addTodo).
		Case(DoneTodo).Then(doneTodo).
		Case(DeleteTodo).Then(deleteTodo).
		Case(ExportTodo).ThenErr(nil, ErrExportDisabled).
		ElseErr(nil, ErrUnsupportedCommand)

	if err != nil {
		fmt.Println("dispatch failed:", err)
		return
	}

	handler()
}

func addTodo() {
	fmt.Println("added todo")
}

func doneTodo() {
	fmt.Println("marked todo as done")
}

func deleteTodo() {
	fmt.Println("deleted todo")
}
