package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestLLDPNeighbors(t *testing.T) {
	// Fixture captured live from lldpd in an LXD container (two nodes on a
	// bridge with LLDP forwarding enabled).
	ns, err := LLDPNeighbors(testutil.Fixture(t, "lldp/neighbors-json0.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(ns) != 1 {
		t.Fatalf("neighbors = %d, want 1", len(ns))
	}
	n := ns[0]
	if n.LocalPort != "eth1" {
		t.Errorf("local port = %q, want eth1", n.LocalPort)
	}
	if n.Chassis != "lldp2.lxd" {
		t.Errorf("chassis = %q, want lldp2.lxd", n.Chassis)
	}
	if n.PortID != "00:16:3e:9b:4d:63" || n.PortDescr != "eth1" || n.TTL != "4" {
		t.Errorf("port = %+v", n)
	}
	if len(n.MgmtIPs) != 2 || n.MgmtIPs[0] != "10.235.37.207" {
		t.Errorf("mgmt-ips = %v", n.MgmtIPs)
	}
	// Only enabled capabilities are kept (Router was enabled).
	found := false
	for _, c := range n.Capabilities {
		if c == "Router" {
			found = true
		}
	}
	if !found {
		t.Errorf("capabilities = %v, want Router", n.Capabilities)
	}
}

func TestLLDPNeighborsEmpty(t *testing.T) {
	ns, err := LLDPNeighbors([]byte(`{"lldp":[{"interface":[]}]}`))
	if err != nil || len(ns) != 0 {
		t.Errorf("ns=%v err=%v", ns, err)
	}
}
