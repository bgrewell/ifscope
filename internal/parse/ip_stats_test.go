package parse

import "testing"

func TestIPLinkStats(t *testing.T) {
	data := []byte(`[
	  {"ifname":"eth0","stats64":{"rx":{"bytes":1048576,"packets":10,"errors":1,"dropped":2},"tx":{"bytes":2048,"packets":5,"errors":0,"dropped":0,"collisions":3}}},
	  {"ifname":"nostats"}
	]`)
	stats, err := IPLinkStats(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(stats) != 1 {
		t.Fatalf("stats = %d, want 1 (interface without stats64 skipped)", len(stats))
	}
	s := stats[0]
	if s.Name != "eth0" || s.RxBytes != 1048576 || s.RxErrors != 1 || s.RxDropped != 2 {
		t.Errorf("rx = %+v", s)
	}
	if s.TxBytes != 2048 || s.TxPackets != 5 || s.Collisions != 3 {
		t.Errorf("tx = %+v", s)
	}
}
