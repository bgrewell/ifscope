package parse

import "testing"

func TestBridgeFDB(t *testing.T) {
	data := []byte(`[
	  {"mac":"33:33:00:00:00:01","ifname":"eth0","flags":["self"],"state":"permanent"},
	  {"mac":"aa:bb:cc:dd:ee:ff","ifname":"eth1","master":"br0","vlan":100,"state":""}
	]`)
	entries, err := BridgeFDB(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("entries = %d, want 2", len(entries))
	}
	if entries[0].MAC != "33:33:00:00:00:01" || entries[0].Dev != "eth0" || entries[0].State != "permanent" {
		t.Errorf("e0 = %+v", entries[0])
	}
	e1 := entries[1]
	if e1.Master != "br0" || e1.VLAN == nil || *e1.VLAN != 100 {
		t.Errorf("e1 = %+v", e1)
	}
}
