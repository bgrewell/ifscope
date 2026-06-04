package collect

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/sysfs"
)

func TestBridgesCollect(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{"br0", "eth0", "eth1", "lo"}

	// br0 is a bridge with STP on and two members.
	fs.Dirs["/sys/class/net/br0/bridge"] = []string{"stp_state"}
	fs.Files["/sys/class/net/br0/bridge/stp_state"] = "1\n"
	fs.Dirs["/sys/class/net/br0/brif"] = []string{"eth0", "eth1"}
	fs.Files["/sys/class/net/br0/operstate"] = "up\n"
	fs.Files["/sys/class/net/eth0/operstate"] = "up\n"
	fs.Files["/sys/class/net/eth1/operstate"] = "down\n"

	bridges, _ := NewBridges(fs).Collect()
	if len(bridges) != 1 {
		t.Fatalf("bridges = %d, want 1", len(bridges))
	}
	b := bridges[0]
	if b.Name != "br0" || b.State != "UP" || !b.STP {
		t.Errorf("bridge = %+v", b)
	}
	if len(b.Members) != 2 || b.Members[0].Name != "eth0" || b.Members[0].State != "UP" {
		t.Errorf("members = %+v", b.Members)
	}
	if b.Members[1].State != "DOWN" {
		t.Errorf("eth1 state = %q, want DOWN", b.Members[1].State)
	}
}

func TestBridgesCollectNone(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Dirs["/sys/class/net"] = []string{"eth0", "lo"}
	fs.Files["/sys/class/net/eth0/operstate"] = "up\n"
	if bridges, _ := NewBridges(fs).Collect(); len(bridges) != 0 {
		t.Errorf("bridges = %d, want 0", len(bridges))
	}
}
