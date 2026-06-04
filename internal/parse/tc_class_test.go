package parse

import "testing"

func TestTCClasses(t *testing.T) {
	// Captured from `tc -json class show` for an htb setup (rate/ceil in bytes/s).
	data := []byte(`[
	  {"class":"htb","handle":"1:10","parent":"1:1","prio":0,"rate":3750000,"ceil":10000000},
	  {"class":"htb","handle":"1:1","root":true,"rate":12500000,"ceil":12500000}
	]`)
	cs, err := TCClasses(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) != 2 {
		t.Fatalf("classes = %d, want 2", len(cs))
	}
	// 3750000 bytes/s * 8 = 30 Mbit/s.
	if cs[0].Handle != "1:10" || cs[0].Rate != 30_000_000 || cs[0].Ceil != 80_000_000 {
		t.Errorf("class0 = %+v", cs[0])
	}
	if !cs[1].Root || cs[1].Rate != 100_000_000 {
		t.Errorf("root class = %+v", cs[1])
	}
}

func TestTCClassesText(t *testing.T) {
	// Older tc emits text for `class show`; must error (caller degrades).
	if _, err := TCClasses([]byte("class htb 1:10 parent 1:1 rate 30Mbit")); err == nil {
		t.Fatal("expected error for non-JSON tc class output")
	}
}
