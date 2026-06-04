package parse

import "testing"

func TestNetnsList(t *testing.T) {
	data := []byte("myns (id: 0)\nother\nthird (id: 5)\n")
	ns := NetnsList(data)
	if len(ns) != 3 {
		t.Fatalf("namespaces = %d, want 3", len(ns))
	}
	if ns[0].Name != "myns" || ns[0].ID == nil || *ns[0].ID != 0 {
		t.Errorf("ns0 = %+v", ns[0])
	}
	if ns[1].Name != "other" || ns[1].ID != nil {
		t.Errorf("ns1 = %+v (id should be nil)", ns[1])
	}
	if ns[2].ID == nil || *ns[2].ID != 5 {
		t.Errorf("ns2 id = %v, want 5", ns[2].ID)
	}
}

func TestNetnsListEmpty(t *testing.T) {
	if ns := NetnsList(nil); len(ns) != 0 {
		t.Errorf("empty input = %v, want none", ns)
	}
}
