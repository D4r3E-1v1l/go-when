package when

import "fmt"

// AnyMatcherRoot is the initial state of a matcher for non-comparable input
// values. It intentionally has no Case method.
type AnyMatcherRoot[T any, R any] struct {
	value T
}

// AnyMatcher is the ready state of a matcher for non-comparable input values.
// It intentionally has no Case method.
type AnyMatcher[T any, R any] struct {
	value   T
	matched bool
	result  R
}

// AnyMatcherCondition is the pending condition state of an AnyMatcher.
type AnyMatcherCondition[T any, R any] struct {
	matcher AnyMatcher[T, R]
	match   func(T) bool
}

func MatchAny[T any, R any](value T) AnyMatcherRoot[T, R] {
	return AnyMatcherRoot[T, R]{value: value}
}

func MatchAnyAs[R any, T any](value T) AnyMatcherRoot[T, R] {
	return AnyMatcherRoot[T, R]{value: value}
}

func (m AnyMatcherRoot[T, R]) ready() AnyMatcher[T, R] {
	return AnyMatcher[T, R]{value: m.value}
}

func (m AnyMatcherRoot[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) AnyMatcherCondition[T, R] {
	return m.ready().Range(first, rest...)
}

func (m AnyMatcherRoot[T, R]) When(first func(T) bool, rest ...func(T) bool) AnyMatcherCondition[T, R] {
	return m.ready().When(first, rest...)
}

func (m AnyMatcherRoot[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) AnyMatcherCondition[T, R] {
	return m.ready().Pattern(first, rest...)
}

func (m AnyMatcher[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) AnyMatcherCondition[T, R] {
	return AnyMatcherCondition[T, R]{
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

func (m AnyMatcher[T, R]) When(first func(T) bool, rest ...func(T) bool) AnyMatcherCondition[T, R] {
	return AnyMatcherCondition[T, R]{
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

func (m AnyMatcher[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) AnyMatcherCondition[T, R] {
	return AnyMatcherCondition[T, R]{
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

func (c AnyMatcherCondition[T, R]) Then(result R) AnyMatcher[T, R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.value) {
		m.matched = true
		m.result = result
	}
	return m
}

func (c AnyMatcherCondition[T, R]) ThenDo(fn func(T) R) AnyMatcher[T, R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.value) {
		m.matched = true
		m.result = fn(m.value)
	}
	return m
}

func (m AnyMatcher[T, R]) Else(result R) R {
	if m.matched {
		return m.result
	}
	return result
}

func (m AnyMatcher[T, R]) ElseDo(fn func(T) R) R {
	if m.matched {
		return m.result
	}
	return fn(m.value)
}

func (m AnyMatcher[T, R]) Exhaustive() R {
	if m.matched {
		return m.result
	}

	panic(fmt.Errorf("%w: value=%v", ErrNoMatch, m.value))
}
