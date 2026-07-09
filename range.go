package when

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidRangeExp = errors.New("when: invalid RangeExp; use when.Range, when.Closed, when.Open, when.LeftOpen, when.From, when.After, when.To, or when.Until")
	ErrInvalidRange    = errors.New("when: invalid range")
)

// RangeExp is a typed range matcher.
//
// The zero value of RangeExp is invalid.
// Use Range, Closed, Open, LeftOpen, From, After, To, or Until to create one.
type RangeExp[T any] struct {
	match func(T) bool
}

func (r RangeExp[T]) Match(v T) bool {
	if r.match == nil {
		panic(ErrInvalidRangeExp)
	}
	return r.match(v)
}

// Range returns a half-open range: [low, high).
func Range[T Numeric](low, high T) RangeExp[T] {
	validateLowLTHigh(low, high)
	return RangeExp[T]{match: func(v T) bool {
		return v >= low && v < high
	}}
}

// Closed returns a closed range: [low, high].
func Closed[T Numeric](low, high T) RangeExp[T] {
	validateLowLEHigh(low, high)
	return RangeExp[T]{match: func(v T) bool {
		return v >= low && v <= high
	}}
}

// Open returns an open range: (low, high).
func Open[T Numeric](low, high T) RangeExp[T] {
	validateLowLTHigh(low, high)
	return RangeExp[T]{match: func(v T) bool {
		return v > low && v < high
	}}
}

// LeftOpen returns a left-open, right-closed range: (low, high].
func LeftOpen[T Numeric](low, high T) RangeExp[T] {
	validateLowLTHigh(low, high)
	return RangeExp[T]{match: func(v T) bool {
		return v > low && v <= high
	}}
}

// From returns [low, +∞).
func From[T Numeric](low T) RangeExp[T] {
	return RangeExp[T]{match: func(v T) bool {
		return v >= low
	}}
}

// After returns (low, +∞).
func After[T Numeric](low T) RangeExp[T] {
	return RangeExp[T]{match: func(v T) bool {
		return v > low
	}}
}

// To returns (-∞, high].
func To[T Numeric](high T) RangeExp[T] {
	return RangeExp[T]{match: func(v T) bool {
		return v <= high
	}}
}

// Until returns (-∞, high).
func Until[T Numeric](high T) RangeExp[T] {
	return RangeExp[T]{match: func(v T) bool {
		return v < high
	}}
}

func validateLowLTHigh[T Numeric](low, high T) {
	if low >= high {
		panic(fmt.Errorf("%w: low=%v high=%v; low must be < high", ErrInvalidRange, low, high))
	}
}

func validateLowLEHigh[T Numeric](low, high T) {
	if low > high {
		panic(fmt.Errorf("%w: low=%v high=%v; low must be <= high", ErrInvalidRange, low, high))
	}
}
