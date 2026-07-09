package when

import (
	"errors"
	"math"
	"testing"
)

func TestRangeHalfOpen(t *testing.T) {
	exp := Range(0, 10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below low", v: -1, want: false},
		{name: "at low", v: 0, want: true},
		{name: "inside", v: 5, want: true},
		{name: "at high", v: 10, want: false},
		{name: "above high", v: 11, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("Range(0, 10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestClosedRange(t *testing.T) {
	exp := Closed(0, 10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below low", v: -1, want: false},
		{name: "at low", v: 0, want: true},
		{name: "inside", v: 5, want: true},
		{name: "at high", v: 10, want: true},
		{name: "above high", v: 11, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("Closed(0, 10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestOpenRange(t *testing.T) {
	exp := Open(0, 10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below low", v: -1, want: false},
		{name: "at low", v: 0, want: false},
		{name: "inside", v: 5, want: true},
		{name: "at high", v: 10, want: false},
		{name: "above high", v: 11, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("Open(0, 10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestLeftOpenRange(t *testing.T) {
	exp := LeftOpen(0, 10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below low", v: -1, want: false},
		{name: "at low", v: 0, want: false},
		{name: "inside", v: 5, want: true},
		{name: "at high", v: 10, want: true},
		{name: "above high", v: 11, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("LeftOpen(0, 10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestFromRange(t *testing.T) {
	exp := From(10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below low", v: 9, want: false},
		{name: "at low", v: 10, want: true},
		{name: "above low", v: 11, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("From(10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestAfterRange(t *testing.T) {
	exp := After(10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below low", v: 9, want: false},
		{name: "at low", v: 10, want: false},
		{name: "above low", v: 11, want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("After(10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestToRange(t *testing.T) {
	exp := To(10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below high", v: 9, want: true},
		{name: "at high", v: 10, want: true},
		{name: "above high", v: 11, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("To(10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestUntilRange(t *testing.T) {
	exp := Until(10)

	tests := []struct {
		name string
		v    int
		want bool
	}{
		{name: "below high", v: 9, want: true},
		{name: "at high", v: 10, want: false},
		{name: "above high", v: 11, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("Until(10).Match(%d) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestRangeWithEqualBounds(t *testing.T) {
	t.Run("Range equal bounds panics", func(t *testing.T) {
		assertPanicsIs(t, func() {
			_ = Range(1, 1)
		}, ErrInvalidRange)
	})

	t.Run("Closed equal bounds matches single value", func(t *testing.T) {
		exp := Closed(1, 1)

		if !exp.Match(1) {
			t.Fatal("Closed(1, 1).Match(1) = false, want true")
		}

		if exp.Match(0) {
			t.Fatal("Closed(1, 1).Match(0) = true, want false")
		}

		if exp.Match(2) {
			t.Fatal("Closed(1, 1).Match(2) = true, want false")
		}
	})

	t.Run("Open equal bounds panics", func(t *testing.T) {
		assertPanicsIs(t, func() {
			_ = Open(1, 1)
		}, ErrInvalidRange)
	})

	t.Run("LeftOpen equal bounds panics", func(t *testing.T) {
		assertPanicsIs(t, func() {
			_ = LeftOpen(1, 1)
		}, ErrInvalidRange)
	})
}

func TestRangePanicsWhenLowGreaterThanHigh(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{name: "Range", fn: func() { _ = Range(10, 0) }},
		{name: "Closed", fn: func() { _ = Closed(10, 0) }},
		{name: "Open", fn: func() { _ = Open(10, 0) }},
		{name: "LeftOpen", fn: func() { _ = LeftOpen(10, 0) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPanicsIs(t, tt.fn, ErrInvalidRange)
		})
	}
}

func TestZeroRangeExpPanics(t *testing.T) {
	assertPanicsIs(t, func() {
		var exp RangeExp[int]
		_ = exp.Match(1)
	}, ErrInvalidRangeExp)
}

func TestFloatRange(t *testing.T) {
	exp := Range(0.5, 0.8)

	tests := []struct {
		name string
		v    float64
		want bool
	}{
		{name: "below low", v: 0.49, want: false},
		{name: "at low", v: 0.5, want: true},
		{name: "inside", v: 0.7, want: true},
		{name: "at high", v: 0.8, want: false},
		{name: "above high", v: 0.81, want: false},
		{name: "nan", v: math.NaN(), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := exp.Match(tt.v); got != tt.want {
				t.Fatalf("Range(0.5, 0.8).Match(%v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestCustomNumericTypeRange(t *testing.T) {
	type Score int

	exp := Range(Score(60), Score(90))

	if !exp.Match(Score(80)) {
		t.Fatal("Range(Score(60), Score(90)).Match(Score(80)) = false, want true")
	}

	if exp.Match(Score(90)) {
		t.Fatal("Range(Score(60), Score(90)).Match(Score(90)) = true, want false")
	}
}

func TestRangeWorksWithMatcher(t *testing.T) {
	got := MatchAs[string](503).
		Case(200).Then("ok").
		Range(Range(500, 600)).Then("server_error").
		Else("unknown")

	if got != "server_error" {
		t.Fatalf("got %q, want %q", got, "server_error")
	}
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	fn()
}

func assertPanicsIs(t *testing.T, fn func(), target error) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic")
		}

		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic value type = %T, want error", recovered)
		}

		if !errors.Is(err, target) {
			t.Fatalf("panic error = %v, want errors.Is(..., %v)", err, target)
		}
	}()

	fn()
}
