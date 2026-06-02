package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestPingAvgMillis(t *testing.T) {
	avg, ok := PingAvgMillis(testutil.Fixture(t, "ping/ok.txt"))
	if !ok {
		t.Fatal("expected rtt to be found")
	}
	if avg != "6.154" {
		t.Errorf("avg = %q, want 6.154", avg)
	}
}

func TestPingAvgMillisLoss(t *testing.T) {
	if _, ok := PingAvgMillis(testutil.Fixture(t, "ping/loss.txt")); ok {
		t.Error("expected no rtt on 100% loss")
	}
}
