# Design Notes

`go-when` explores a narrow space between Go's native `switch` and full pattern matching.

It is intentionally limited.

## Not a better switch

`go-when` is not intended to be a better `switch`.

Go's `if` and `switch` are often the clearest way to express procedural logic.

`go-when` focuses on flat, result-oriented decision mapping.

## Decision table helper

A good `go-when` chain usually looks like a small decision table:

```text
condition -> result
condition -> result
condition -> result
fallback  -> result
```

This structure is useful for:

- value mapping
- handler dispatch
- error mapping
- range classification
- mixed condition matching

## Keep control flow explicit

`go-when` does not try to model `return`, `break`, or `continue`.

Those are language-level control-flow features and are better expressed with native Go.

## Why no custom Result type

Go already has a strong error-return convention:

```go
value, err := f()
```

`go-when` follows this convention with `WithErr()`.

It does not introduce a custom `Result[T, E]` type.

## Why explicit terminals

A matcher chain should make fallback or exhaustive intent explicit.

Use `Else` when fallback behavior is expected.

Use `Exhaustive` when all valid cases are expected to be covered.

This helps readers understand what happens when no branch matches.

## First match wins

Like Go `switch` and Rust `match`, `go-when` evaluates branches in order.

The first matched condition wins.

This keeps the model simple and predictable.
