package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/run"
)

func TestRulesCollect(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: `[{"priority":0,"src":"all","table":"local"},{"priority":100,"src":"192.168.9.0","srclen":24,"table":"50"}]`},
		"ip", "-json", "rule",
	)
	rules, warnings := NewRules(fake).Collect(context.Background())
	if len(warnings) != 0 {
		t.Fatalf("warnings: %v", warnings)
	}
	if len(rules) != 2 {
		t.Fatalf("rules = %d, want 2", len(rules))
	}
	if rules[1].From != "192.168.9.0/24" || rules[1].Table != "50" {
		t.Errorf("rule = %+v", rules[1])
	}
}

func TestRulesCollectMissingIP(t *testing.T) {
	_, warnings := NewRules(run.NewFake()).Collect(context.Background())
	if len(warnings) != 1 || !warnings[0].Fatal {
		t.Fatalf("warnings = %v, want one fatal", warnings)
	}
}
