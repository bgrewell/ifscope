package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestBridgeMDB(t *testing.T) {
	// Fixture captured from a live container (static + auto entries).
	entries, err := BridgeMDB(testutil.Fixture(t, "bridge/mdb.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(entries))
	}
	var found bool
	for _, e := range entries {
		if e.Bridge == "br0" && e.Port == "veth0" && e.Group == "239.1.1.1" && e.State == "permanent" {
			found = true
		}
	}
	if !found {
		t.Errorf("expected permanent 239.1.1.1 on br0/veth0; got %+v", entries)
	}
}

func TestBridgeMDBEmpty(t *testing.T) {
	entries, err := BridgeMDB([]byte(`[{"mdb":[],"router":{}}]`))
	if err != nil || len(entries) != 0 {
		t.Errorf("entries=%v err=%v", entries, err)
	}
}
