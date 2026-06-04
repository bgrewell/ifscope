package parse

import "testing"

func TestMulticastGroups(t *testing.T) {
	data := []byte(`[
	  {"ifindex":2,"ifname":"eth0","maddr":[
	    {"link":"33:33:00:00:00:01"},
	    {"family":"inet","address":"224.0.0.1"},
	    {"family":"inet6","address":"ff02::1"}
	  ]}
	]`)
	groups, err := MulticastGroups(data)
	if err != nil {
		t.Fatal(err)
	}
	// The link-only (MAC) entry is skipped; two IP groups remain.
	if len(groups) != 2 {
		t.Fatalf("groups = %d, want 2", len(groups))
	}
	if groups[0].Interface != "eth0" || groups[0].Family != "inet" || groups[0].Address != "224.0.0.1" {
		t.Errorf("group0 = %+v", groups[0])
	}
	if groups[1].Family != "inet6" {
		t.Errorf("group1 family = %q", groups[1].Family)
	}
}
