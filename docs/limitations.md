# Limitations

`go-when` is intentionally limited.

It is not a full pattern matching library and does not replace Go's `if` or `switch`.

## Not for complex control flow

Prefer native Go when branches need:

- `return`
- `break`
- `continue`
- many statements
- local state mutation
- complex side effects

Native Go is often clearer for procedural logic.

## Not for deep structural pattern matching

`go-when` does not try to provide Rust-style or TypeScript-style deep pattern matching.

Prefer extracting a decision value first, then matching it.

```go
state := order.State

action := when.MatchAs[Action](state).
	Case(OrderPaid).Then(PackOrder).
	Else(Noop)
```

## First match wins

Matchers are evaluated in order.

Put specific conditions before broad conditions.

```go
grade := when.MatchAs[string](score).
	Case(100).Then("perfect").
	Range(when.Range(90, 100)).Then("excellent").
	Else("unknown")
```

## `When` is flexible but not analyzable

`When(func(T) bool)` can express custom business conditions, but static tools cannot reliably analyze arbitrary predicate functions.

Use `Case` and `Range` when the condition is simple and analyzable.

Use `When` when the condition is domain-specific.

## Prefer native switch for short condition switches

This is already clear Go:

```go
switch {
case x.Valid && x.Int64 == 50:
	fmt.Println("Got 50")
case x.Valid:
	fmt.Println("Matched", x.Int64)
default:
	fmt.Println("default", x)
}
```

`go-when` is more useful for flat, result-oriented mappings.

## Not for ultra-hot paths

Prefer native `if` or `switch` in extremely hot loops or low-level performance-sensitive paths.

`go-when` is intended for business decision mapping, not for replacing low-level branching in tight loops.
