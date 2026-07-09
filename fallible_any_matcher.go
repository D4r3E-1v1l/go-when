package when

import "fmt"

// FallibleAnyMatcherRoot is the initial state of a fallible matcher for any
// input value. It intentionally has no Case method.
//
// Use MatchAny or MatchAnyAs followed by WithErr to enter this mode:
//
//	result, err := when.MatchAnyAs[Decision](req).
//		WithErr().
//		When(isValid).Then(Allow).
//		ElseErr(Deny, ErrInvalidRequest)
type FallibleAnyMatcherRoot[T any, R any] struct {
	value T
}

// FallibleAnyMatcher is the ready state of a fallible matcher for any input
// value. It intentionally has no Case method.
//
// Matching is first-match-wins. Once a branch matches, following branches are
// skipped, but the chain can continue until a terminal method returns (R, error).
type FallibleAnyMatcher[T any, R any] struct {
	value   T
	matched bool
	result  R
	err     error
}

// FallibleAnyMatcherCondition is the pending condition state of a
// FallibleAnyMatcher.
//
// It must be completed with Then, ThenErr, ThenDo, or ThenDoE before the chain
// can continue.
type FallibleAnyMatcherCondition[T any, R any] struct {
	matcher FallibleAnyMatcher[T, R]
	match   func(T) bool
}

// WithErr converts an any-value matcher root into a fallible any-value matcher
// root.
//
// WithErr is intentionally only available on the root state. This keeps the
// chain shape clear:
//
//	MatchAnyAs[R](value).WithErr().When(...).ThenErr(...).ElseErr(...)
func (m AnyMatcherRoot[T, R]) WithErr() FallibleAnyMatcherRoot[T, R] {
	return FallibleAnyMatcherRoot[T, R]{value: m.value}
}

func (m FallibleAnyMatcherRoot[T, R]) ready() FallibleAnyMatcher[T, R] {
	return FallibleAnyMatcher[T, R]{value: m.value}
}

func (m FallibleAnyMatcherRoot[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) FallibleAnyMatcherCondition[T, R] {
	return m.ready().Range(first, rest...)
}

func (m FallibleAnyMatcherRoot[T, R]) When(first func(T) bool, rest ...func(T) bool) FallibleAnyMatcherCondition[T, R] {
	return m.ready().When(first, rest...)
}

func (m FallibleAnyMatcherRoot[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) FallibleAnyMatcherCondition[T, R] {
	return m.ready().Pattern(first, rest...)
}

func (m FallibleAnyMatcher[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) FallibleAnyMatcherCondition[T, R] {
	return FallibleAnyMatcherCondition[T, R]{
		matcher: m,
		match: func(v T) bool {
			if first.Match(v) {
				return true
			}
			for _, exp := range rest {
				if exp.Match(v) {
					return true
				}
			}
			return false
		},
	}
}

func (m FallibleAnyMatcher[T, R]) When(first func(T) bool, rest ...func(T) bool) FallibleAnyMatcherCondition[T, R] {
	return FallibleAnyMatcherCondition[T, R]{
		matcher: m,
		match: func(v T) bool {
			if first(v) {
				return true
			}
			for _, pred := range rest {
				if pred(v) {
					return true
				}
			}
			return false
		},
	}
}

func (m FallibleAnyMatcher[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) FallibleAnyMatcherCondition[T, R] {
	return FallibleAnyMatcherCondition[T, R]{
		matcher: m,
		match: func(v T) bool {
			if first.Match(v) {
				return true
			}
			for _, pattern := range rest {
				if pattern.Match(v) {
					return true
				}
			}
			return false
		},
	}
}

// Then completes a matched branch with result and nil error.
func (c FallibleAnyMatcherCondition[T, R]) Then(result R) FallibleAnyMatcher[T, R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.value) {
		m.matched = true
		m.result = result
		m.err = nil
	}
	return m
}

// ThenErr completes a matched branch with result and err.
func (c FallibleAnyMatcherCondition[T, R]) ThenErr(result R, err error) FallibleAnyMatcher[T, R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.value) {
		m.matched = true
		m.result = result
		m.err = err
	}
	return m
}

// ThenDo lazily computes result for a matched branch and returns nil error.
func (c FallibleAnyMatcherCondition[T, R]) ThenDo(fn func(T) R) FallibleAnyMatcher[T, R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.value) {
		m.matched = true
		m.result = fn(m.value)
		m.err = nil
	}
	return m
}

// ThenDoE lazily computes result and error for a matched branch.
func (c FallibleAnyMatcherCondition[T, R]) ThenDoE(fn func(T) (R, error)) FallibleAnyMatcher[T, R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.value) {
		m.matched = true
		m.result, m.err = fn(m.value)
	}
	return m
}

// Else returns result and nil error when no previous branch matched.
func (m FallibleAnyMatcher[T, R]) Else(result R) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return result, nil
}

// ElseErr returns result and err when no previous branch matched.
func (m FallibleAnyMatcher[T, R]) ElseErr(result R, err error) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return result, err
}

// ElseDo lazily computes result and returns nil error when no previous branch
// matched.
func (m FallibleAnyMatcher[T, R]) ElseDo(fn func(T) R) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return fn(m.value), nil
}

// ElseDoE lazily computes result and error when no previous branch matched.
func (m FallibleAnyMatcher[T, R]) ElseDoE(fn func(T) (R, error)) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return fn(m.value)
}

// Exhaustive returns the matched result and error, or returns ErrNoMatch when no
// branch matched.
func (m FallibleAnyMatcher[T, R]) Exhaustive() (R, error) {
	if m.matched {
		return m.result, m.err
	}

	var zero R
	return zero, fmt.Errorf("%w: value=%v", ErrNoMatch, m.value)
}
