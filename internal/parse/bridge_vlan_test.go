package parse

import "testing"

func TestBridgeVLANs(t *testing.T) {
	data := []byte(`[
	  {"ifname":"lxdbr0","vlans":[{"vlan":1,"flags":["PVID","Egress Untagged"]}]},
	  {"ifname":"eth0","vlans":[{"vlan":10,"flags":[]},{"vlan":100,"vlanEnd":200,"flags":[]}]}
	]`)
	vlans, err := BridgeVLANs(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(vlans) != 3 {
		t.Fatalf("entries = %d, want 3", len(vlans))
	}
	if vlans[0].Port != "lxdbr0" || vlans[0].VLAN != 1 || len(vlans[0].Flags) != 2 {
		t.Errorf("entry0 = %+v", vlans[0])
	}
	// VLAN range preserved.
	if vlans[2].VLAN != 100 || vlans[2].VLANEnd != 200 {
		t.Errorf("range entry = %+v", vlans[2])
	}
}
