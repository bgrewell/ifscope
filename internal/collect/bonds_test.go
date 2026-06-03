package collect

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/sysfs"
)

// bondFS builds a fake sysfs with two bonds: an up 802.3ad bond with two up
// members, and a down active-backup bond with two down members. Plain NICs and
// a loopback are present to confirm only bonds are detected.
func bondFS() *sysfs.Fake {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{
		"bond.les3", "bond.les4", "enp67s0f0", "enp67s0f1", "enp67s0f2", "enp67s0f3", "lo",
	}

	set := func(name, key, val string) {
		fs.Files["/sys/class/net/"+name+"/"+key] = val + "\n"
	}

	// bond.les3 — LACP, up.
	set("bond.les3", "operstate", "up")
	set("bond.les3", "bonding/slaves", "enp67s0f2 enp67s0f3")
	set("bond.les3", "bonding/mode", "802.3ad 4")
	set("bond.les3", "bonding/active_slave", "enp67s0f2")
	set("enp67s0f2", "operstate", "up")
	set("enp67s0f3", "operstate", "up")

	// bond.les4 — active-backup, down.
	set("bond.les4", "operstate", "down")
	set("bond.les4", "bonding/slaves", "enp67s0f1 enp67s0f0")
	set("bond.les4", "bonding/mode", "active-backup 1")
	set("enp67s0f0", "operstate", "down")
	set("enp67s0f1", "operstate", "down")

	// A plain NIC (no bonding dir) must be ignored.
	set("lo", "operstate", "unknown")
	return fs
}

func TestBondsCollect(t *testing.T) {
	bonds, _ := NewBonds(bondFS()).Collect()
	if len(bonds) != 2 {
		t.Fatalf("bonds = %d, want 2", len(bonds))
	}

	les3 := bonds[0] // sorted by name: bond.les3 < bond.les4
	if les3.Name != "bond.les3" {
		t.Fatalf("first bond = %q", les3.Name)
	}
	if les3.Mode != "802.3ad (LACP)" {
		t.Errorf("mode = %q, want 802.3ad (LACP)", les3.Mode)
	}
	if les3.State != "UP" {
		t.Errorf("state = %q, want UP", les3.State)
	}
	if les3.ActiveSlave != "enp67s0f2" {
		t.Errorf("active = %q", les3.ActiveSlave)
	}
	if len(les3.Members) != 2 || les3.Members[0].Name != "enp67s0f2" || les3.Members[0].State != "UP" {
		t.Errorf("members = %+v", les3.Members)
	}

	les4 := bonds[1]
	if les4.Mode != "active-backup" || les4.State != "DOWN" {
		t.Errorf("les4 = %+v", les4)
	}
	// Members are sorted by name regardless of the slaves-file order.
	if les4.Members[0].Name != "enp67s0f0" || les4.Members[0].State != "DOWN" {
		t.Errorf("les4 members = %+v", les4.Members)
	}
	if les4.ActiveSlave != "" {
		t.Errorf("les4 active = %q, want empty", les4.ActiveSlave)
	}
}

func TestBondsCollectNone(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{"eth0", "lo"}
	fs.Files["/sys/class/net/eth0/operstate"] = "up\n"
	bonds, _ := NewBonds(fs).Collect()
	if len(bonds) != 0 {
		t.Errorf("bonds = %d, want 0", len(bonds))
	}
}
