package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestOVS(t *testing.T) {
	ovs, err := OVS(
		testutil.Fixture(t, "ovs/bridge.json"),
		testutil.Fixture(t, "ovs/port.json"),
		testutil.Fixture(t, "ovs/interface.json"),
	)
	if err != nil {
		t.Fatalf("OVS: %v", err)
	}

	if len(ovs.Bridges) != 1 {
		t.Fatalf("bridges = %d, want 1", len(ovs.Bridges))
	}
	br := ovs.Bridges[0]
	if len(br.Ports) != 5 {
		t.Errorf("bridge %q ports = %d, want 5", br.Name, len(br.Ports))
	}

	// Every port should resolve a name, back-reference its bridge, and the
	// access VLAN port should carry its tag.
	byName := map[string]model.OVSPort{}
	for _, p := range ovs.Ports {
		byName[p.Name] = p
		if p.Bridge != br.Name {
			t.Errorf("port %q bridge = %q, want %q", p.Name, p.Bridge, br.Name)
		}
		if len(p.Interfaces) == 0 {
			t.Errorf("port %q resolved no interfaces", p.Name)
		}
	}

	vlan30, ok := byName["vlan30"]
	if !ok {
		t.Fatal("expected a vlan30 port")
	}
	if vlan30.Tag == nil || *vlan30.Tag != 30 {
		t.Errorf("vlan30 tag = %v, want 30", vlan30.Tag)
	}
}

func TestOVSSetUnwrapping(t *testing.T) {
	// A port with a single interface uses the unwrapped ["uuid",...] shape.
	bridge := []byte(`{"headings":["name","ports"],"data":[["br0",["uuid","p1"]]]}`)
	port := []byte(`{"headings":["_uuid","name","interfaces","tag","trunks","vlan_mode","bond_mode","lacp"],"data":[[["uuid","p1"],"eth0",["uuid","i1"],["set",[]],["set",[]],["set",[]],["set",[]],["set",[]]]]}`)
	iface := []byte(`{"headings":["_uuid","name","type","ofport","external_ids","options"],"data":[[["uuid","i1"],"eth0","",1,["map",[]],["map",[]]]]}`)

	ovs, err := OVS(bridge, port, iface)
	if err != nil {
		t.Fatal(err)
	}
	if len(ovs.Ports) != 1 || ovs.Ports[0].Name != "eth0" {
		t.Fatalf("ports = %+v", ovs.Ports)
	}
	if got := ovs.Ports[0].Interfaces; len(got) != 1 || got[0] != "eth0" {
		t.Errorf("interfaces = %v, want [eth0]", got)
	}
	if ovs.Bridges[0].Ports[0] != "eth0" {
		t.Errorf("bridge port = %q, want eth0", ovs.Bridges[0].Ports[0])
	}
}

func TestOVSTrunksAndMap(t *testing.T) {
	port := []byte(`{"headings":["_uuid","name","interfaces","tag","trunks","vlan_mode","bond_mode","lacp"],"data":[[["uuid","p1"],"trunk0",["uuid","i1"],["set",[]],["set",[100,200,300]],"trunk",["set",[]],["set",[]]]]}`)
	iface := []byte(`{"headings":["_uuid","name","type","ofport","external_ids","options"],"data":[[["uuid","i1"],"trunk0","internal",2,["map",[["foo","bar"]]],["map",[]]]]}`)
	bridge := []byte(`{"headings":["name","ports"],"data":[["br0",["uuid","p1"]]]}`)

	ovs, err := OVS(bridge, port, iface)
	if err != nil {
		t.Fatal(err)
	}
	p := ovs.Ports[0]
	if len(p.Trunks) != 3 || p.Trunks[0] != 100 {
		t.Errorf("trunks = %v, want [100 200 300]", p.Trunks)
	}
	if p.VLANMode != "trunk" {
		t.Errorf("vlan_mode = %q", p.VLANMode)
	}
	if ovs.Interfaces[0].ExternalIDs["foo"] != "bar" {
		t.Errorf("external_ids = %v", ovs.Interfaces[0].ExternalIDs)
	}
}
