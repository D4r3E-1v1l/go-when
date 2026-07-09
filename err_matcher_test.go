package when

import (
	"errors"
	"fmt"
	"testing"
)

var (
	errNotFound      = errors.New("not found")
	errAlreadyExists = errors.New("already exists")
)

type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}

func TestErrNil(t *testing.T) {
	got := Err[string](nil).
		Nil().Then("ok").
		Else("failed")

	if got != "ok" {
		t.Fatalf("got %q, want %q", got, "ok")
	}
}

func TestErrNilDoIsLazy(t *testing.T) {
	called := false

	got := Err[string](errNotFound).
		Nil().ThenDo(func(error) string {
		called = true
		return "ok"
	}).
		Else("failed")

	if called {
		t.Fatal("Nil handler was called for non-nil error")
	}
	if got != "failed" {
		t.Fatalf("got %q, want %q", got, "failed")
	}
}

func TestErrNotNil(t *testing.T) {
	got := Err[string](errNotFound).
		NotNil().Then("failed").
		Else("ok")

	if got != "failed" {
		t.Fatalf("got %q, want %q", got, "failed")
	}
}

func TestErrNotNilDoIsLazy(t *testing.T) {
	called := false

	got := Err[string](nil).
		NotNil().ThenDo(func(err error) string {
		called = true
		return err.Error()
	}).
		Else("ok")

	if called {
		t.Fatal("NotNil handler was called for nil error")
	}
	if got != "ok" {
		t.Fatalf("got %q, want %q", got, "ok")
	}
}

func TestErrIs(t *testing.T) {
	err := fmt.Errorf("wrap: %w", errNotFound)

	got := Err[string](err).
		Is(errAlreadyExists).Then("conflict").
		Is(errNotFound).Then("not_found").
		Else("internal")

	if got != "not_found" {
		t.Fatalf("got %q, want %q", got, "not_found")
	}
}

func TestErrIsSupportsMultipleTargets(t *testing.T) {
	err := fmt.Errorf("wrap: %w", errNotFound)

	got := Err[string](err).
		Is(errAlreadyExists, errNotFound).Then("known_error").
		Else("internal")

	if got != "known_error" {
		t.Fatalf("got %q, want %q", got, "known_error")
	}
}

func TestErrIsFirstMatchWins(t *testing.T) {
	err := fmt.Errorf("wrap: %w", errNotFound)

	got := Err[string](err).
		Is(errNotFound).Then("first").
		Is(errNotFound).Then("second").
		Else("unknown")

	if got != "first" {
		t.Fatalf("got %q, want %q", got, "first")
	}
}

func TestErrIsDoIsLazy(t *testing.T) {
	called := false

	got := Err[string](errNotFound).
		Is(errAlreadyExists).ThenDo(func(err error) string {
		called = true
		return "conflict"
	}).
		Else("other")

	if called {
		t.Fatal("Is handler was called for unmatched error")
	}
	if got != "other" {
		t.Fatalf("got %q, want %q", got, "other")
	}
}

func TestErrAs(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &validationError{msg: "invalid sandbox id"})

	var ve *validationError

	got := Err[string](err).
		As(&ve).Then("bad_request").
		Else("internal")

	if got != "bad_request" {
		t.Fatalf("got %q, want %q", got, "bad_request")
	}
	if ve == nil {
		t.Fatal("validation error target was not populated")
	}
	if ve.msg != "invalid sandbox id" {
		t.Fatalf("ve.msg = %q, want %q", ve.msg, "invalid sandbox id")
	}
}

func TestErrAsDo(t *testing.T) {
	err := fmt.Errorf("wrap: %w", &validationError{msg: "invalid cidr"})

	var ve *validationError

	got := Err[string](err).
		As(&ve).ThenDo(func(error) string {
		return "bad_request: " + ve.msg
	}).
		Else("internal")

	want := "bad_request: invalid cidr"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestErrContains(t *testing.T) {
	err := errors.New("dial tcp: connection refused")

	got := Err[string](err).
		Contains("timeout").Then("gateway_timeout").
		Contains("connection refused").Then("unavailable").
		Else("internal")

	if got != "unavailable" {
		t.Fatalf("got %q, want %q", got, "unavailable")
	}
}

func TestErrContainsSupportsMultipleSubstrings(t *testing.T) {
	err := errors.New("dial tcp: connection refused")

	got := Err[string](err).
		Contains("timeout", "connection refused").Then("unavailable").
		Else("internal")

	if got != "unavailable" {
		t.Fatalf("got %q, want %q", got, "unavailable")
	}
}

func TestErrContainsWithNilErrorDoesNotMatch(t *testing.T) {
	got := Err[string](nil).
		Contains("anything").Then("matched").
		Else("not_matched")

	if got != "not_matched" {
		t.Fatalf("got %q, want %q", got, "not_matched")
	}
}

func TestErrWhen(t *testing.T) {
	err := errors.New("quota exceeded")

	got := Err[string](err).
		When(func(err error) bool {
			return err != nil && err.Error() == "quota exceeded"
		}).Then("too_many_requests").
		Else("internal")

	if got != "too_many_requests" {
		t.Fatalf("got %q, want %q", got, "too_many_requests")
	}
}

func TestErrWhenDoIsLazy(t *testing.T) {
	called := false

	got := Err[string](errors.New("x")).
		When(func(err error) bool {
			return false
		}).ThenDo(func(err error) string {
		called = true
		return "matched"
	}).
		Else("unmatched")

	if called {
		t.Fatal("When handler was called for unmatched predicate")
	}
	if got != "unmatched" {
		t.Fatalf("got %q, want %q", got, "unmatched")
	}
}

type timeoutPattern struct{}

func (timeoutPattern) Match(err error) bool {
	return err != nil && err.Error() == "timeout"
}

func TestErrPattern(t *testing.T) {
	got := Err[string](errors.New("timeout")).
		Pattern(timeoutPattern{}).Then("gateway_timeout").
		Else("internal")

	if got != "gateway_timeout" {
		t.Fatalf("got %q, want %q", got, "gateway_timeout")
	}
}

func TestErrPatternDo(t *testing.T) {
	got := Err[string](errors.New("timeout")).
		Pattern(timeoutPattern{}).ThenDo(func(err error) string {
		return "matched: " + err.Error()
	}).
		Else("internal")

	want := "matched: timeout"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestErrElseDo(t *testing.T) {
	err := errors.New("unknown")

	got := Err[string](err).
		Is(errNotFound).Then("not_found").
		ElseDo(func(err error) string {
			return "internal: " + err.Error()
		})

	want := "internal: unknown"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}
