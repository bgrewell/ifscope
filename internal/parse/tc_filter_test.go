package parse

import "testing"

func TestTCFilters(t *testing.T) {
	// Captured shape: header + hash-table entries (no flowid) plus the rule.
	data := []byte(`[
	  {"parent":"1:","protocol":"ip","pref":1,"kind":"u32","chain":0},
	  {"parent":"1:","protocol":"ip","pref":1,"kind":"u32","options":{"fh":"800:","ht_divisor":1}},
	  {"parent":"1:","protocol":"ip","pref":1,"kind":"u32","options":{"fh":"800::800","flowid":"1:10","match":{"value":"50","mask":"ffff","off":20}}},
	  {"parent":"1:","protocol":"ip","pref":2,"kind":"u32","options":{"fh":"801::800","flowid":"1:20"}}
	]`)
	fs, err := TCFilters(data)
	if err != nil {
		t.Fatal(err)
	}
	// Only the two rule entries with a flowid are kept.
	if len(fs) != 2 {
		t.Fatalf("filters = %d, want 2", len(fs))
	}
	if fs[0].Pref != 1 || fs[0].Kind != "u32" || fs[0].FlowID != "1:10" {
		t.Errorf("filter0 = %+v", fs[0])
	}
	if fs[1].FlowID != "1:20" {
		t.Errorf("filter1 flowid = %q", fs[1].FlowID)
	}
}

func TestTCFiltersText(t *testing.T) {
	if _, err := TCFilters([]byte("filter parent 1: protocol ip pref 1 u32")); err == nil {
		t.Fatal("expected error for non-JSON tc filter output")
	}
}
