package when

import (
	"errors"
	"fmt"
	"strings"
)

// ErrPredicate is a function-based error matcher.
type ErrPredicate func(error) bool

func (p ErrPredicate) Match(err error) bool {
	return p(err)
}

// ErrPattern is implemented by custom error matchers.
type ErrPattern interface {
	Match(error) bool
}

// ErrMatcherRoot is the initial state of an error matcher.
type ErrMatcherRoot[R any] struct {
	err error
}

// ErrMatcher is the ready state of an error matcher.
//
// Matching is first-match-wins. Once a branch matches, following branches are
// skipped, but the chain can continue until a terminal method returns R.
type ErrMatcher[R any] struct {
	err     error
	matched bool
	result  R
}

// ErrMatcherCondition is the pending condition state of an ErrMatcher.
// It must be completed with Then or ThenDo before the chain can continue.
type ErrMatcherCondition[R any] struct {
	matcher ErrMatcher[R]
	match   func(error) bool
}

// Err creates an error matcher root.
func Err[R any](err error) ErrMatcherRoot[R] {
	return ErrMatcherRoot[R]{err: err}
}

func (m ErrMatcherRoot[R]) ready() ErrMatcher[R] {
	return ErrMatcher[R]{err: m.err}
}

// Nil matches a nil error.
func (m ErrMatcherRoot[R]) Nil() ErrMatcherCondition[R] {
	return m.ready().Nil()
}

// NotNil matches any non-nil error.
func (m ErrMatcherRoot[R]) NotNil() ErrMatcherCondition[R] {
	return m.ready().NotNil()
}

// Is matches errors.Is(err, target) for one or more targets.
func (m ErrMatcherRoot[R]) Is(first error, rest ...error) ErrMatcherCondition[R] {
	return m.ready().Is(first, rest...)
}

// As matches errors.As(err, target).
//
// target must be a non-nil pointer accepted by errors.As. For typed access,
// declare the target variable outside and read it in ThenDo:
//
//	var ve *ValidationError
//	resp := when.Err[Resp](err).
//	    As(&ve).ThenDo(func(error) Resp {
//	        return badRequest(ve)
//	    }).
//	    Else(internal)
func (m ErrMatcherRoot[R]) As(target any) ErrMatcherCondition[R] {
	return m.ready().As(target)
}

// Contains matches strings.Contains(err.Error(), substr) for one or more
// substrings.
//
// This is useful for fallback/error-message based matching, but errors.Is,
// errors.As, or domain error codes are usually better for stable APIs.
func (m ErrMatcherRoot[R]) Contains(first string, rest ...string) ErrMatcherCondition[R] {
	return m.ready().Contains(first, rest...)
}

// When matches one or more custom predicates.
func (m ErrMatcherRoot[R]) When(first func(error) bool, rest ...func(error) bool) ErrMatcherCondition[R] {
	return m.ready().When(first, rest...)
}

// Pattern matches one or more reusable error matchers.
func (m ErrMatcherRoot[R]) Pattern(first ErrPattern, rest ...ErrPattern) ErrMatcherCondition[R] {
	return m.ready().Pattern(first, rest...)
}

func (m ErrMatcher[R]) Nil() ErrMatcherCondition[R] {
	return ErrMatcherCondition[R]{
		matcher: m,
		match: func(err error) bool {
			return err == nil
		},
	}
}

func (m ErrMatcher[R]) NotNil() ErrMatcherCondition[R] {
	return ErrMatcherCondition[R]{
		matcher: m,
		match: func(err error) bool {
			return err != nil
		},
	}
}

func (m ErrMatcher[R]) Is(first error, rest ...error) ErrMatcherCondition[R] {
	return ErrMatcherCondition[R]{
		matcher: m,
		match: func(err error) bool {
			if errors.Is(err, first) {
				return true
			}
			for _, target := range rest {
				if errors.Is(err, target) {
					return true
				}
			}
			return false
		},
	}
}

func (m ErrMatcher[R]) As(target any) ErrMatcherCondition[R] {
	return ErrMatcherCondition[R]{
		matcher: m,
		match: func(err error) bool {
			return errors.As(err, target)
		},
	}
}

func (m ErrMatcher[R]) Contains(first string, rest ...string) ErrMatcherCondition[R] {
	return ErrMatcherCondition[R]{
		matcher: m,
		match: func(err error) bool {
			if err == nil {
				return false
			}
			msg := err.Error()
			if strings.Contains(msg, first) {
				return true
			}
			for _, substr := range rest {
				if strings.Contains(msg, substr) {
					return true
				}
			}
			return false
		},
	}
}

func (m ErrMatcher[R]) When(first func(error) bool, rest ...func(error) bool) ErrMatcherCondition[R] {
	return ErrMatcherCondition[R]{
		matcher: m,
		match: func(err error) bool {
			if first(err) {
				return true
			}
			for _, pred := range rest {
				if pred(err) {
					return true
				}
			}
			return false
		},
	}
}

func (m ErrMatcher[R]) Pattern(first ErrPattern, rest ...ErrPattern) ErrMatcherCondition[R] {
	return ErrMatcherCondition[R]{
		matcher: m,
		match: func(err error) bool {
			if first.Match(err) {
				return true
			}
			for _, pattern := range rest {
				if pattern.Match(err) {
					return true
				}
			}
			return false
		},
	}
}

func (c ErrMatcherCondition[R]) Then(result R) ErrMatcher[R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.err) {
		m.matched = true
		m.result = result
	}
	return m
}

func (c ErrMatcherCondition[R]) ThenDo(fn func(error) R) ErrMatcher[R] {
	m := c.matcher
	if m.matched {
		return m
	}
	if c.match(m.err) {
		m.matched = true
		m.result = fn(m.err)
	}
	return m
}

// Else returns result when no previous branch matched.
func (m ErrMatcher[R]) Else(result R) R {
	if m.matched {
		return m.result
	}
	return result
}

// ElseDo lazily computes the result when no previous branch matched.
func (m ErrMatcher[R]) ElseDo(fn func(error) R) R {
	if m.matched {
		return m.result
	}
	return fn(m.err)
}

// Exhaustive returns the matched result or panics when no branch matched.
func (m ErrMatcher[R]) Exhaustive() R {
	if m.matched {
		return m.result
	}

	panic(fmt.Errorf("%w: err=%v", ErrNoMatch, m.err))
}
