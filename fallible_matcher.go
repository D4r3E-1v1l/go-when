package when

import "fmt"

// FallibleMatcherRoot is the initial state of a comparable fallible matcher.
//
// Use Match or MatchAs followed by WithErr to enter this mode:
//
//	result, err := when.MatchAs[Resp](code).
//		WithErr().
//		Case(200).Then(okResp).
//		Case(400).ThenErr(badRequestResp, ErrBadRequest).
//		ElseErr(internalResp, ErrUnknown)
//
// A fallible matcher maps an input value to (R, error).
type FallibleMatcherRoot[T comparable, R any] struct {
	value T
}

// FallibleMatcher is the ready state of a comparable fallible matcher.
//
// Matching is first-match-wins. Once a branch matches, following branches are
// skipped, but the chain can continue until a terminal method returns (R, error).
type FallibleMatcher[T comparable, R any] struct {
	value   T
	matched bool
	result  R
	err     error
}

// FallibleMatcherCondition is the pending condition state of a comparable
// fallible matcher.
//
// It must be completed with Then, ThenErr, ThenDo, or ThenDoE before the chain
// can continue.
type FallibleMatcherCondition[T comparable, R any] struct {
	matcher FallibleMatcher[T, R]
	match   func(T) bool
}

// WithErr converts a comparable matcher root into a fallible matcher root.
//
// WithErr is intentionally only available on the root state. This keeps the
// chain shape clear:
//
//	MatchAs[R](value).WithErr().Case(...).ThenErr(...).ElseErr(...)
func (m MatcherRoot[T, R]) WithErr() FallibleMatcherRoot[T, R] {
	return FallibleMatcherRoot[T, R]{value: m.value}
}

func (m FallibleMatcherRoot[T, R]) ready() FallibleMatcher[T, R] {
	return FallibleMatcher[T, R]{value: m.value}
}

func (m FallibleMatcherRoot[T, R]) Case(first T, rest ...T) FallibleMatcherCondition[T, R] {
	return m.ready().Case(first, rest...)
}

func (m FallibleMatcherRoot[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) FallibleMatcherCondition[T, R] {
	return m.ready().Range(first, rest...)
}

func (m FallibleMatcherRoot[T, R]) When(first func(T) bool, rest ...func(T) bool) FallibleMatcherCondition[T, R] {
	return m.ready().When(first, rest...)
}

func (m FallibleMatcherRoot[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) FallibleMatcherCondition[T, R] {
	return m.ready().Pattern(first, rest...)
}

func (m FallibleMatcher[T, R]) Case(first T, rest ...T) FallibleMatcherCondition[T, R] {
	return FallibleMatcherCondition[T, R]{
		matcher: m,
		match: func(v T) bool {
			if v == first {
				return true
			}
			for _, expected := range rest {
				if v == expected {
					return true
				}
			}
			return false
		},
	}
}

func (m FallibleMatcher[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) FallibleMatcherCondition[T, R] {
	return FallibleMatcherCondition[T, R]{
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

func (m FallibleMatcher[T, R]) When(first func(T) bool, rest ...func(T) bool) FallibleMatcherCondition[T, R] {
	return FallibleMatcherCondition[T, R]{
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

func (m FallibleMatcher[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) FallibleMatcherCondition[T, R] {
	return FallibleMatcherCondition[T, R]{
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
func (c FallibleMatcherCondition[T, R]) Then(result R) FallibleMatcher[T, R] {
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
func (c FallibleMatcherCondition[T, R]) ThenErr(result R, err error) FallibleMatcher[T, R] {
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
func (c FallibleMatcherCondition[T, R]) ThenDo(fn func(T) R) FallibleMatcher[T, R] {
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
func (c FallibleMatcherCondition[T, R]) ThenDoE(fn func(T) (R, error)) FallibleMatcher[T, R] {
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
func (m FallibleMatcher[T, R]) Else(result R) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return result, nil
}

// ElseErr returns result and err when no previous branch matched.
func (m FallibleMatcher[T, R]) ElseErr(result R, err error) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return result, err
}

// ElseDo lazily computes result and returns nil error when no previous branch
// matched.
func (m FallibleMatcher[T, R]) ElseDo(fn func(T) R) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return fn(m.value), nil
}

// ElseDoE lazily computes result and error when no previous branch matched.
func (m FallibleMatcher[T, R]) ElseDoE(fn func(T) (R, error)) (R, error) {
	if m.matched {
		return m.result, m.err
	}
	return fn(m.value)
}

// Exhaustive returns the matched result and error, or returns ErrNoMatch when no
// branch matched.
func (m FallibleMatcher[T, R]) Exhaustive() (R, error) {
	if m.matched {
		return m.result, m.err
	}

	var zero R
	return zero, fmt.Errorf("%w: value=%v", ErrNoMatch, m.value)
}
