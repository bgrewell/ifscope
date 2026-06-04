package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestTunnels(t *testing.T) {
	// Fixture captured from a live container (vxlan/gre/geneve), including the
	// "fan-map" token that breaks `ip -json`.
	tuns := Tunnels(testutil.Fixture(t, "ip/link-tunnels.txt"))
	if len(tuns) != 3 {
		t.Fatalf("tunnels = %d, want 3", len(tuns))
	}
	by := map[string]struct {
		typ, local, remote, ttl string
		vni, port               int
		hasVNI, hasPort         bool
	}{}
	for _, x := range tuns {
		e := by[x.Name]
		e.typ, e.local, e.remote, e.ttl = x.Type, x.Local, x.Remote, x.TTL
		if x.VNI != nil {
			e.vni, e.hasVNI = *x.VNI, true
		}
		if x.Port != nil {
			e.port, e.hasPort = *x.Port, true
		}
		by[x.Name] = e
	}

	vx := by["vxlan100"]
	if vx.typ != "vxlan" || !vx.hasVNI || vx.vni != 100 || vx.remote != "192.0.2.2" || !vx.hasPort || vx.port != 4789 {
		t.Errorf("vxlan100 = %+v", vx)
	}
	gre := by["gre1"]
	if gre.typ != "gre" || gre.local != "10.235.37.149" || gre.remote != "192.0.2.3" || gre.ttl != "inherit" {
		t.Errorf("gre1 = %+v", gre)
	}
	gnv := by["gnv200"]
	if gnv.typ != "geneve" || !gnv.hasVNI || gnv.vni != 200 || gnv.port != 6081 {
		t.Errorf("gnv200 = %+v", gnv)
	}
}

func TestTunnelsNone(t *testing.T) {
	data := []byte("1: lo: <LOOPBACK,UP> mtu 65536\n    link/loopback 00:00:00:00:00:00\n")
	if got := Tunnels(data); len(got) != 0 {
		t.Errorf("tunnels = %v, want none", got)
	}
}
