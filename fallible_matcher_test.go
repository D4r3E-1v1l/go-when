package when

import (
	"errors"
	"testing"
)

func TestMatcherWithErrThenReturnsNilError(t *testing.T) {
	got, err := MatchAs[string](200).
		WithErr().
		Case(200).Then("ok").
		ElseErr("unknown", errors.New("unknown status"))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != "ok" {
		t.Fatalf("got %q", got)
	}
}

func TestMatcherWithErrThenErr(t *testing.T) {
	wantErr := errors.New("bad request")

	got, err := MatchAs[string](400).
		WithErr().
		Case(200).Then("ok").
		Case(400).ThenErr("bad_request", wantErr).
		ElseErr("unknown", errors.New("unknown status"))

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if got != "bad_request" {
		t.Fatalf("got %q", got)
	}
}

func TestMatcherWithErrElseErr(t *testing.T) {
	wantErr := errors.New("unknown status")

	got, err := MatchAs[string](418).
		WithErr().
		Case(200).Then("ok").
		ElseErr("unknown", wantErr)

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if got != "unknown" {
		t.Fatalf("got %q", got)
	}
}

func TestMatcherWithErrThenDoEIsLazy(t *testing.T) {
	called := false

	got, err := MatchAs[string](200).
		WithErr().
		Case(400).ThenDoE(func(v int) (string, error) {
		called = true
		return "bad_request", errors.New("bad request")
	}).
		Else("ok")

	if called {
		t.Fatal("unmatched handler was called")
	}
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != "ok" {
		t.Fatalf("got %q", got)
	}
}

func TestMatcherWithErrExhaustiveNoMatch(t *testing.T) {
	_, err := MatchAs[string](404).
		WithErr().
		Case(200).Then("ok").
		Exhaustive()

	if !errors.Is(err, ErrNoMatch) {
		t.Fatalf("expected ErrNoMatch, got %v", err)
	}
}
