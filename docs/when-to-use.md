# When to use go-when

Use `go-when` when the decision can be expressed as a flat mapping:

- value -> value
- value -> handler
- value -> `(value, error)`
- error -> value or handler
- numeric range -> category
- mixed `Case` / `Range` / `When` / `Pattern` conditions

`go-when` is not designed to replace Go's `if` or `switch`. It is a helper for narrow cases where a typed decision chain is easier to scan than repeated assignment or repeated dispatch code.

## Good fit: value to value

```go
label := when.MatchAs[string](status).
	Case("pending").Then("Waiting for payment").
	Case("paid").Then("Payment received").
	Case("shipped").Then("On the way").
	Else("Unknown status")
```

This is a good fit because each branch is simply:

```text
condition -> result
```

## Good fit: value to handler

```go
type Handler func()

handler := when.MatchAs[Handler](state).
	Case(OrderCreated).Then(requestPayment).
	Case(OrderPaid).Then(packOrder).
	Case(OrderPacked).Then(shipOrder).
	Case(OrderCancelled).Then(cancelOrder).
	Exhaustive()

handler()
```

The matcher only selects the handler. The actual execution remains explicit.

## Good fit: range classification

```go
grade := when.MatchAs[string](score).
	Case(100).Then("perfect").
	Range(when.Range(90, 100)).Then("excellent").
	Range(when.Range(60, 90)).Then("passed").
	Else("failed")
```

This is a good fit because the result is a category selected from exact values and ranges.

## Good fit: custom predicate helpers

```go
decision := when.MatchAnyAs[string](form).
	When(isInvalidEmail).Then("invalid_email").
	When(needsParentApproval).Then("needs_parent_approval").
	When(isAllowedPlan).Then("allow").
	Else("manual_review")
```

Prefer named helper functions over large inline anonymous functions.

## Prefer native if/switch

Prefer native `if` or `switch` when the logic is procedural, multi-step, or control-flow-heavy.

```go
func decide(form SignupForm) Decision {
	if invalidEmail(form) {
		return InvalidEmail
	}

	if needsParentApproval(form) {
		return NeedsParentApproval
	}

	if allowedPlan(form) {
		return Allow
	}

	return ManualReview
}
```

This is clear Go code. `go-when` is not intended to replace this style.

## Rule of thumb

Use `go-when` when each branch is mostly:

```text
condition -> result
```

Use native Go when each branch contains:

```text
multiple statements
local state changes
logging/metrics
return/break/continue
complex side effects
```
