# Mixed Condition Matching

`go-when` allows exact value cases, numeric ranges, custom predicates, and reusable patterns to be combined in one typed decision chain.

This is one of the main reasons to use `go-when` instead of a plain value switch or condition switch.

## Example

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

This chain combines:

- exact value matching with `Case`
- numeric range matching with `Range`
- custom predicate matching with `When`
- explicit fallback with `Else`

## First match wins

Conditions are evaluated in order.

The first matched condition wins.

Put more specific conditions before broader conditions.

Good:

```go
grade := when.MatchAs[string](score).
	Case(100).Then("perfect").
	Range(when.Range(90, 100)).Then("excellent").
	Else("unknown")
```

Bad:

```go
grade := when.MatchAs[string](score).
	Range(when.Range(90, 101)).Then("excellent").
	Case(100).Then("perfect").
	Else("unknown")
```

In the second example, `100` is already covered by the range.

## Helper functions for custom predicates

Prefer named helper functions when using `When`.

```go
decision := when.MatchAnyAs[string](form).
	When(isInvalidEmail).Then("invalid_email").
	When(needsParentApproval).Then("needs_parent_approval").
	When(isAllowedPlan).Then("allow").
	Else("manual_review")
```

```go
func needsParentApproval(form SignupForm) bool {
	return form.Age < 18 && form.Plan == "pro"
}
```

This keeps the matcher chain readable and makes the condition logic testable.

## Case and Range are more analyzable

`Case` and `Range` are more structured than arbitrary predicates.

They are better candidates for static checks such as:

- duplicate exact values
- range overlap
- unreachable numeric conditions

`When(func(T) bool)` is intentionally flexible, but arbitrary predicate functions are difficult to analyze statically.
