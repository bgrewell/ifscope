package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestIPAddresses(t *testing.T) {
	ifaces, err := IPAddresses(testutil.Fixture(t, "ip/address.json"))
	if err != nil {
		t.Fatalf("IPAddresses: %v", err)
	}
	if len(ifaces) != 7 {
		t.Fatalf("got %d interfaces, want 7", len(ifaces))
	}

	by := make(map[string]model.Interface, len(ifaces))
	for _, i := range ifaces {
		by[i.Name] = i
	}

	t.Run("classification", func(t *testing.T) {
		want := map[string]model.InterfaceType{
			"lo":       model.TypeLoopback,
			"eth0":     model.TypePhysical,
			"eth1":     model.TypePhysical,
			"eth0.100": model.TypeVLAN,
			"br0":      model.TypeBridge,
			"veth1234": model.TypeVeth,
			"vlan2001": model.TypeOVS,
		}
		for name, typ := range want {
			if got := by[name].Type; got != typ {
				t.Errorf("%s type = %q, want %q", name, got, typ)
			}
		}
	})

	t.Run("physical bus and state", func(t *testing.T) {
		eth0 := by["eth0"]
		if eth0.Bus != "0000:33:00.0" {
			t.Errorf("eth0 bus = %q, want 0000:33:00.0", eth0.Bus)
		}
		if eth0.State != "UP" {
			t.Errorf("eth0 state = %q, want UP", eth0.State)
		}
		if eth1 := by["eth1"]; eth1.State != "DOWN" {
			t.Errorf("eth1 state = %q, want DOWN", eth1.State)
		}
	})

	t.Run("addresses keep both families", func(t *testing.T) {
		eth0 := by["eth0"]
		if len(eth0.Addresses) != 3 {
			t.Fatalf("eth0 addresses = %d, want 3", len(eth0.Addresses))
		}
		if v4 := eth0.IPv4(); len(v4) != 2 {
			t.Errorf("eth0 IPv4 = %d, want 2", len(v4))
		}
		if !eth0.Addresses[1].Secondary {
			t.Errorf("eth0 second address should be secondary")
		}
	})

	t.Run("vlan parent and tag", func(t *testing.T) {
		v := by["eth0.100"]
		if v.LinkParent != "eth0" {
			t.Errorf("vlan parent = %q, want eth0", v.LinkParent)
		}
		if v.VLANID != 100 {
			t.Errorf("vlan id = %d, want 100", v.VLANID)
		}
	})

	t.Run("master membership", func(t *testing.T) {
		if veth := by["veth1234"]; veth.LinkParent != "br0" {
			t.Errorf("veth master = %q, want br0", veth.LinkParent)
		}
	})
}

func TestIPAddressesVirtualKinds(t *testing.T) {
	// Virtual link kinds without a dedicated constant must report the kernel
	// kind verbatim, not be mislabeled physical.
	data := []byte(`[
	  {"ifindex":2,"ifname":"macvlan0","link_type":"ether","address":"aa:aa:aa:aa:aa:aa",
	   "linkinfo":{"info_kind":"macvlan"}},
	  {"ifindex":3,"ifname":"vx0","link_type":"ether","address":"bb:bb:bb:bb:bb:bb",
	   "linkinfo":{"info_kind":"vxlan"}},
	  {"ifindex":4,"ifname":"wg0","link_type":"none",
	   "linkinfo":{"info_kind":"wireguard"}},
	  {"ifindex":5,"ifname":"weird0","link_type":"ether","address":"cc:cc:cc:cc:cc:cc"}
	]`)
	ifaces, err := IPAddresses(data)
	if err != nil {
		t.Fatal(err)
	}
	by := map[string]model.InterfaceType{}
	for _, i := range ifaces {
		by[i.Name] = i.Type
	}
	if by["macvlan0"] != "macvlan" {
		t.Errorf("macvlan0 type = %q, want macvlan", by["macvlan0"])
	}
	if by["vx0"] != "vxlan" {
		t.Errorf("vx0 type = %q, want vxlan", by["vx0"])
	}
	if by["wg0"] != "wireguard" {
		t.Errorf("wg0 type = %q, want wireguard", by["wg0"])
	}
	// ether, no linkinfo, no device backing → undetermined, not assumed physical.
	if by["weird0"] != model.TypeUnknown {
		t.Errorf("weird0 type = %q, want unknown", by["weird0"])
	}
}

func TestIPAddressesInvalidJSON(t *testing.T) {
	if _, err := IPAddresses([]byte("not json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestIPAddressesRealHostFixture(t *testing.T) {
	// The captured host fixture must parse without error and yield interfaces.
	ifaces, err := IPAddresses(testutil.Fixture(t, "ip/address-host.json"))
	if err != nil {
		t.Fatalf("IPAddresses(host): %v", err)
	}
	if len(ifaces) == 0 {
		t.Fatal("host fixture produced no interfaces")
	}
}
