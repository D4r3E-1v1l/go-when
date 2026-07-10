# Fallible Matcher

Use `WithErr()` when the decision should return `(R, error)`.

This follows Go's native error style and avoids introducing a custom `Result[T, E]` type.

## Basic example

```go
handler, err := when.MatchAs[Handler](command).
	WithErr().
	Case(AddTodo).Then(addTodo).
	Case(DeleteTodo).Then(deleteTodo).
	Case(ExportTodo).ThenErr(nil, ErrExportDisabled).
	ElseErr(nil, ErrUnsupportedCommand)
```

The returned values are normal Go values:

```go
if err != nil {
	return err
}

handler()
```

## Then vs ThenErr

Use `Then` when the branch succeeds.

```go
Case(AddTodo).Then(addTodo)
```

Use `ThenErr` when the branch should return an error.

```go
Case(ExportTodo).ThenErr(nil, ErrExportDisabled)
```

## ElseErr

Use `ElseErr` for unsupported or unexpected values.

```go
ElseErr(nil, ErrUnsupportedCommand)
```

## ThenDoE

Use `ThenDoE` when the result needs to be computed and the computation may fail.

```go
result, err := when.MatchAs[Result](command).
	WithErr().
	Case(ImportData).ThenDoE(importData).
	ElseErr(Result{}, ErrUnsupportedCommand)
```

The function passed to `ThenDoE` should have this shape:

```go
func(T) (R, error)
```

## Keep the return model simple

`go-when` intentionally uses Go's native return style:

```go
value, err := ...
```

It does not introduce a custom `Result[T, E]` type.
