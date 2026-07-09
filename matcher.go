package when

import (
	"errors"
	"fmt"
)

var ErrNoMatch = errors.New("when: no match")

// Predicate is a function-based matcher.
type Predicate[T any] func(T) bool

func (p Predicate[T]) Match(v T) bool {
	return p(v)
}

// Pattern is implemented by custom matchers.
type Pattern[T any] interface {
	Match(T) bool
}

// MatcherRoot is the initial state of a comparable matcher.
//
// A root can only open a condition. Once a condition is completed with Then or
// ThenDo, it becomes a Matcher and can either open another condition or finish
// with Else, ElseDo, or Exhaustive.
type MatcherRoot[T comparable, R any] struct {
	value T
}

// Matcher is the ready state of a comparable matcher.
//
// Matching is first-match-wins. Once a branch matches, following branches are
// skipped, but the chain can continue until a terminal method returns R.
type Matcher[T comparable, R any] struct {
	value   T
	matched bool
	result  R
}

// MatcherCondition is the pending condition state of a comparable matcher.
//
// It must be completed with Then or ThenDo before the chain can continue.
type MatcherCondition[T comparable, R any] struct {
	matcher Matcher[T, R]
	match   func(T) bool
}

// Match creates a matcher root for comparable values.
func Match[T comparable, R any](value T) MatcherRoot[T, R] {
	return MatcherRoot[T, R]{value: value}
}

// MatchAs lets callers specify only the result type:
//
//	when.MatchAs[string](code)
func MatchAs[R any, T comparable](value T) MatcherRoot[T, R] {
	return MatcherRoot[T, R]{value: value}
}

func (m MatcherRoot[T, R]) ready() Matcher[T, R] {
	return Matcher[T, R]{value: m.value}
}

func (m MatcherRoot[T, R]) Case(first T, rest ...T) MatcherCondition[T, R] {
	return m.ready().Case(first, rest...)
}

func (m MatcherRoot[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) MatcherCondition[T, R] {
	return m.ready().Range(first, rest...)
}

func (m MatcherRoot[T, R]) When(first func(T) bool, rest ...func(T) bool) MatcherCondition[T, R] {
	return m.ready().When(first, rest...)
}

func (m MatcherRoot[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) MatcherCondition[T, R] {
	return m.ready().Pattern(first, rest...)
}

func (m Matcher[T, R]) Case(first T, rest ...T) MatcherCondition[T, R] {
	return MatcherCondition[T, R]{
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

func (m Matcher[T, R]) Range(first RangeExp[T], rest ...RangeExp[T]) MatcherCondition[T, R] {
	return MatcherCondition[T, R]{
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

func (m Matcher[T, R]) When(first func(T) bool, rest ...func(T) bool) MatcherCondition[T, R] {
	return MatcherCondition[T, R]{
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

func (m Matcher[T, R]) Pattern(first Pattern[T], rest ...Pattern[T]) MatcherCondition[T, R] {
	return MatcherCondition[T, R]{
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

func (c MatcherCondition[T, R]) Then(result R) Matcher[T, R] {
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

func (c MatcherCondition[T, R]) ThenDo(fn func(T) R) Matcher[T, R] {
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

func (m Matcher[T, R]) Else(result R) R {
	if m.matched {
		return m.result
	}
	return result
}

func (m Matcher[T, R]) ElseDo(fn func(T) R) R {
	if m.matched {
		return m.result
	}
	return fn(m.value)
}

func (m Matcher[T, R]) Exhaustive() R {
	if m.matched {
		return m.result
	}

	panic(fmt.Errorf("%w: value=%v", ErrNoMatch, m.value))
}
