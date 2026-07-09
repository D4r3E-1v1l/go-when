package when

import "testing"

func TestMatchAny(t *testing.T) {
	type Req struct {
		Tags []string
	}

	got := MatchAnyAs[string](Req{Tags: []string{"create"}}).
		When(func(req Req) bool {
			return len(req.Tags) > 0 && req.Tags[0] == "create"
		}).Then("create").
		Else("unknown")

	if got != "create" {
		t.Fatalf("got %q", got)
	}
}
