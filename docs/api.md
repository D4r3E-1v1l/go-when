# API Overview

This document gives a high-level overview of the main `go-when` APIs.

## MatchAs

Use `MatchAs` for comparable values.

```go
result := when.MatchAs[string](value).
	Case(A).Then("a").
	Case(B).Then("b").
	Else("unknown")
```

Typical use cases:

- status -> label
- enum -> result
- command -> handler
- state -> action

## MatchAnyAs

Use `MatchAnyAs` for values that are not necessarily comparable, usually structs.

```go
decision := when.MatchAnyAs[string](form).
	When(isInvalidEmail).Then("invalid_email").
	When(needsParentApproval).Then("needs_parent_approval").
	Else("manual_review")
```

`MatchAnyAs` is useful when the decision is based on helper predicates.

## Err

Use `Err` for error mapping.

```go
result := when.Err[APIResult](err).
	Is(ErrNotFound).Then(notFoundResult).
	Is(ErrInvalidInput).Then(badRequestResult).
	Else(internalErrorResult)
```

`Err` is useful for mapping domain errors to result objects, response builders, or handlers.

## WithErr

Use `WithErr()` when the matcher should return `(R, error)`.

```go
handler, err := when.MatchAs[Handler](command).
	WithErr().
	Case(AddTodo).Then(addTodo).
	Case(ExportTodo).ThenErr(nil, ErrExportDisabled).
	ElseErr(nil, ErrUnsupportedCommand)
```

This follows Go's normal error style.

## Conditions

### Case

Exact value matching.

```go
Case("paid").Then("Payment received")
```

`Case` can accept multiple values.

```go
Case("paid", "shipped").Then("active")
```

### Range

Numeric range matching.

```go
Range(when.Range(60, 90)).Then("passed")
Range(when.From(90)).Then("high")
Range(when.Until(60)).Then("low")
```

### When

Custom predicate matching.

```go
When(isInvalidEmail).Then("invalid_email")
```

Prefer named helper functions for readability.

### Pattern

Reusable custom pattern matching.

```go
Pattern(VIPCart{}).Then(20)
```

Use `Pattern` when the same matching logic is useful in multiple places.

## Actions

Normal matcher actions:

- `Then(value R)`
- `ThenDo(func(T) R)`

Fallible matcher actions:

- `Then(value R)`
- `ThenErr(value R, err error)`
- `ThenDo(func(T) R)`
- `ThenDoE(func(T) (R, error))`

## Terminals

Normal matcher terminals:

- `Else(value R)`
- `ElseDo(func(T) R)`
- `Exhaustive()`

Fallible matcher terminals:

- `Else(value R)`
- `ElseErr(value R, err error)`
- `ElseDo(func(T) R)`
- `ElseDoE(func(T) (R, error))`
- `Exhaustive()`

Every completed matcher chain should end with a terminal method.
