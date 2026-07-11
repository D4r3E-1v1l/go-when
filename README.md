# go-when

![Minimum Go Version](https://img.shields.io/badge/min%20go-1.22-blue)
[![LICENSE](https://img.shields.io/github/license/D4r3E-1v1l/go-when.svg)](https://github.com/D4r3E-1v1l/go-when/blob/main/LICENSE)

[中文文档](./README.zh-CN.md)

`go-when` is a typed decision-table helper for Go.

It is designed for flat, result-oriented decision mappings, such as value mapping, handler dispatch, fallible decisions, error mapping, numeric range classification, and mixed condition matching.

It is not a replacement for Go's `if` or `switch`, and it is not a full pattern matching library.

## Install

```bash
go get github.com/D4r3E-1v1l/go-when
```

## Requirements

Go 1.18 or later.

## Quick Example

```go
label := when.MatchAs[string](status).
	Case("pending").Then("Waiting for payment").
	Case("paid").Then("Payment received").
	Case("shipped").Then("On the way").
	Else("Unknown status")
```

## Why

Go's `if` and `switch` are clear and powerful.

`go-when` is useful when a decision is flat, typed, and result-oriented:

- value -> value
- value -> handler
- value -> `(value, error)`
- error -> value or handler
- numeric range -> category
- mixed `Case` / `Range` / `When` / `Pattern` conditions

Prefer native `if` or `switch` for procedural logic, multi-step branches, complex control flow, or branches that need `return`, `break`, or `continue`.

## Mixed Condition Matching

One of the main goals of `go-when` is to let different condition types live in one typed decision chain.

```go
grade := when.MatchAs[string](score).
	Case(100).Then("perfect").
	Range(when.Range(90, 100)).Then("excellent").
	Range(when.Range(60, 90)).Then("passed").
	When(isRetakeAllowed).Then("retake_allowed").
	Else("failed")
```

```go
func isRetakeAllowed(score int) bool {
	return score >= 50 && score < 60
}
```

The first matched condition wins. Put more specific conditions before broader conditions.

## Value to Handler

`go-when` can be used to map a state or command to a handler.

```go
type Handler func()

handler := when.MatchAs[Handler](state).
	Case(OrderCreated).Then(requestPayment).
	Case(OrderPaid).Then(packOrder).
	Case(OrderPacked).Then(shipOrder).
	Case(OrderShipped).Then(sendTracking).
	Case(OrderCancelled).Then(cancelOrder).
	Exhaustive()

handler()
```

This keeps the matcher focused on the decision:

```text
state -> handler
```

The actual execution stays explicit:

```go
handler()
```

## Error Mapping

Use `Err` to map errors to values or handlers.

```go
handler := when.Err[ErrorHandler](err).
	Nil().Then(success).
	Is(ErrUserNotFound).Then(notFound).
	Is(ErrInvalidInput).Then(badRequest).
	Is(ErrNotAuthorized).Then(unauthorized).
	Else(internalError)

result := handler(err)
```

## Fallible Decisions

Use `WithErr()` when a decision should return `(R, error)`.

```go
handler, err := when.MatchAs[Handler](command).
	WithErr().
	Case(AddTodo).Then(addTodo).
	Case(DeleteTodo).Then(deleteTodo).
	Case(ExportTodo).ThenErr(nil, ErrExportDisabled).
	ElseErr(nil, ErrUnsupportedCommand)
```

This follows Go's native error style:

```go
value, err := ...
```

`go-when` does not introduce a custom `Result[T, E]` type.

## Explicit Terminals

Every completed matcher chain must end with a terminal method.

Use `Else` when fallback behavior is expected:

```go
label := when.MatchAs[string](status).
	Case("paid").Then("Payment received").
	Else("Unknown status")
```

Use `Exhaustive` when all valid cases are expected to be covered:

```go
action := when.MatchAs[Action](state).
	Case(OrderCreated).Then(RequestPayment).
	Case(OrderPaid).Then(PackOrder).
	Case(OrderCancelled).Then(Noop).
	Exhaustive()
```

## Examples

Runnable examples are available in [`examples/`](./examples).

```bash
go run ./examples/01_status_label
go run ./examples/02_order_state_handler
go run ./examples/03_score_grade
go run ./examples/04_error_to_result_handler
go run ./examples/05_command_dispatch_with_error
go run ./examples/06_signup_decision
go run ./examples/07_custom_pattern
```

## Documentation

See [`docs/`](./docs).

Recommended reading:

- [When to use go-when](./docs/when-to-use.md)
- [API Overview](./docs/api.md)
- [Mixed Condition Matching](./docs/mixed-condition-matching.md)
- [Error Mapping](./docs/error-mapping.md)
- [Fallible Matcher](./docs/fallible-matcher.md)
- [Limitations](./docs/limitations.md)

## GoLand Plugin

A companion GoLand plugin is available:

```text
https://github.com/D4r3E-1v1l/go-when-goland-plugin
```

The plugin adds inspections for matcher-chain structure and selected semantic checks, such as missing terminals, numeric overlap, unreachable numeric conditions, and enum exhaustive warnings.

The library works without the plugin.

## License

MIT
