package correlate

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
)

func ifaces() []model.Interface {
	return []model.Interface{
		{Name: "eth1", Type: model.TypePhysical, State: "DOWN", Driver: "ice"},
		{Name: "vlan100", Type: model.TypeVLAN, State: "UP", LinkParent: "eth0", VLANID: 100},
		{Name: "eth0", Type: model.TypePhysical, State: "UP", Driver: "mlx5_core"},
		{Name: "vlan50", Type: model.TypeVLAN, State: "UP", LinkParent: "eth0", VLANID: 50},
		{Name: "lo", Type: model.TypeLoopback, State: "UNKNOWN"},
	}
}

func TestPartition(t *testing.T) {
	in, vlans := Partition(ifaces())
	if len(in) != 3 {
		t.Errorf("interfaces = %d, want 3", len(in))
	}
	if len(vlans) != 2 {
		t.Errorf("vlans = %d, want 2", len(vlans))
	}
	for _, v := range vlans {
		if v.Type != model.TypeVLAN {
			t.Errorf("partition put non-VLAN %q in vlans", v.Name)
		}
	}
}

func TestSortInterfaces(t *testing.T) {
	in, _ := Partition(ifaces())
	SortInterfaces(in)
	want := []string{"eth0", "eth1", "lo"} // UP, then DOWN, then UNKNOWN
	for i, name := range want {
		if in[i].Name != name {
			t.Errorf("position %d = %q, want %q", i, in[i].Name, name)
		}
	}
}

func TestSortVLANs(t *testing.T) {
	_, vlans := Partition(ifaces())
	SortVLANs(vlans)
	if vlans[0].VLANID != 50 || vlans[1].VLANID != 100 {
		t.Errorf("VLAN order = [%d %d], want [50 100]", vlans[0].VLANID, vlans[1].VLANID)
	}
}

func TestFilter(t *testing.T) {
	all := ifaces()

	t.Run("up only", func(t *testing.T) {
		got := Filter{Up: true}.Apply(all)
		if len(got) != 3 {
			t.Errorf("up filter = %d, want 3", len(got))
		}
	})

	t.Run("by name", func(t *testing.T) {
		got := Filter{Name: "eth0"}.Apply(all)
		if len(got) != 1 || got[0].Name != "eth0" {
			t.Errorf("name filter = %v, want [eth0]", got)
		}
	})

	t.Run("by driver", func(t *testing.T) {
		got := Filter{Driver: "ice"}.Apply(all)
		if len(got) != 1 || got[0].Name != "eth1" {
			t.Errorf("driver filter = %v, want [eth1]", got)
		}
	})

	t.Run("state case-insensitive", func(t *testing.T) {
		got := Filter{State: "up"}.Apply(all)
		if len(got) != 3 {
			t.Errorf("state filter = %d, want 3", len(got))
		}
	})

	t.Run("physical excludes vlan and loopback", func(t *testing.T) {
		got := Filter{Physical: true}.Apply(all)
		if len(got) != 2 {
			t.Errorf("physical filter = %d, want 2", len(got))
		}
	})

	t.Run("virtual excludes physical", func(t *testing.T) {
		got := Filter{Virtual: true}.Apply(all)
		if len(got) != 3 {
			t.Errorf("virtual filter = %d, want 3", len(got))
		}
	})
}
