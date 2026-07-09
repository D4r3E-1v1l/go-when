package when

import "testing"

func TestCaseFirstMatchWins(t *testing.T) {
	got := MatchAs[string](200).
		Case(200).Then("ok").
		Case(200).Then("duplicate").
		Else("unknown")

	if got != "ok" {
		t.Fatalf("got %q", got)
	}
}

func TestCaseSupportsMultipleExpectedValues(t *testing.T) {
	got := MatchAs[string](201).
		Case(200, 201, 204).Then("success").
		Case(400, 404).Then("client_error").
		Else("unknown")

	if got != "success" {
		t.Fatalf("got %q", got)
	}
}

func TestRange(t *testing.T) {
	got := MatchAs[string](503).
		Case(200).Then("ok").
		Range(Range(500, 600)).Then("server_error").
		Else("unknown")

	if got != "server_error" {
		t.Fatalf("got %q", got)
	}
}

func TestRangeSupportsMultipleRangeExpressions(t *testing.T) {
	got := MatchAs[string](31000).
		Range(
			Range(0, 1024),
			Range(30000, 32768),
			Range(49152, 65536),
		).Then("special").
		Else("normal")

	if got != "special" {
		t.Fatalf("got %q", got)
	}
}

func TestThenDoIsLazy(t *testing.T) {
	called := false

	got := MatchAs[string](1).
		When(func(v int) bool { return v == 2 }).ThenDo(func(v int) string {
		called = true
		return "two"
	}).
		Else("other")

	if called {
		t.Fatal("unmatched handler was called")
	}
	if got != "other" {
		t.Fatalf("got %q", got)
	}
}
