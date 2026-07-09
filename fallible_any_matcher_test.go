package when

import (
	"errors"
	"testing"
)

func TestAnyMatcherWithErr(t *testing.T) {
	type Req struct {
		Admin bool
	}

	got, err := MatchAnyAs[string](Req{Admin: true}).
		WithErr().
		When(func(req Req) bool {
			return req.Admin
		}).Then("allow").
		ElseErr("deny", errors.New("permission denied"))

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got != "allow" {
		t.Fatalf("got %q", got)
	}
}

func TestAnyMatcherWithErrElseErr(t *testing.T) {
	type Req struct {
		Admin bool
	}

	wantErr := errors.New("permission denied")

	got, err := MatchAnyAs[string](Req{Admin: false}).
		WithErr().
		When(func(req Req) bool {
			return req.Admin
		}).Then("allow").
		ElseErr("deny", wantErr)

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
	if got != "deny" {
		t.Fatalf("got %q", got)
	}
}
