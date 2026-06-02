package parse

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/testutil"
)

func TestVFAttrs(t *testing.T) {
	vfs, err := VFAttrs(testutil.Fixture(t, "ip/link-vf.json"))
	if err != nil {
		t.Fatalf("VFAttrs: %v", err)
	}
	if len(vfs) != 2 {
		t.Fatalf("vfs = %d, want 2", len(vfs))
	}

	v0 := vfs[0]
	if v0.MAC != "aa:bb:cc:dd:ee:00" || v0.VLAN != 100 || !v0.SpoofCheck || v0.Trust || v0.LinkState != "auto" {
		t.Errorf("vf0 = %+v", v0)
	}

	v1 := vfs[1]
	if v1.VLAN != 0 || v1.SpoofCheck || !v1.Trust || v1.LinkState != "enable" {
		t.Errorf("vf1 = %+v", v1)
	}
}

func TestVFAttrsNoVFs(t *testing.T) {
	vfs, err := VFAttrs([]byte(`[{"ifname":"eth0","vfinfo_list":[]}]`))
	if err != nil {
		t.Fatal(err)
	}
	if len(vfs) != 0 {
		t.Errorf("vfs = %d, want 0", len(vfs))
	}
}
