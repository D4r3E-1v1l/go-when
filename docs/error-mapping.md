# Error Mapping

Use `Err` when you want to map an error to a value or handler.

This is useful in service and API-style code where domain errors should be converted into a consistent result.

## Error to result

```go
result := when.Err[APIResult](err).
	Nil().Then(APIResult{Status: 200, Message: "success"}).
	Is(ErrUserNotFound).Then(APIResult{Status: 404, Message: "user not found"}).
	Is(ErrInvalidInput).Then(APIResult{Status: 400, Message: "invalid input"}).
	Else(APIResult{Status: 500, Message: "internal error"})
```

## Error to handler

For more realistic code, you may map an error to a handler and then execute the handler.

```go
type ErrorHandler func(error) APIResult

handler := when.Err[ErrorHandler](err).
	Nil().Then(success).
	Is(ErrUserNotFound).Then(notFound).
	Is(ErrInvalidInput).Then(badRequest).
	Is(ErrNotAuthorized).Then(unauthorized).
	Else(internalError)

result := handler(err)
```

This keeps the decision and execution separate:

```text
error -> handler -> result
```

## Matching wrapped errors

Use `Is` for errors that may be wrapped.

```go
result := when.Err[APIResult](err).
	Is(ErrUserNotFound).Then(notFoundResult).
	Else(internalErrorResult)
```

## Matching error types

Use `As` when the behavior depends on a concrete error type.

```go
var validationErr *ValidationError

result := when.Err[APIResult](err).
	As(&validationErr).ThenDo(func(error) APIResult {
		return APIResult{
			Status:  400,
			Message: "invalid " + validationErr.Field,
		}
	}).
	Else(APIResult{Status: 500, Message: "internal error"})
```

## Matching text

Use `Contains` only when matching text is intentional.

```go
result := when.Err[APIResult](err).
	Contains("timeout").Then(timeoutResult).
	Else(internalErrorResult)
```

Prefer `Is` or `As` when the error type or sentinel error is available.
