package collect

import (
	"context"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
	"github.com/bgrewell/ifscope/internal/run"
	"github.com/bgrewell/ifscope/internal/sysfs"
	"github.com/bgrewell/ifscope/internal/testutil"
)

// sriovFS builds a fake sysfs tree for one PF (pf0, bus 25:00.0) with two VFs:
// VF0 bound to iavf with netdev vf0, VF1 bound to vfio-pci with no netdev.
func sriovFS() *sysfs.Fake {
	fs := sysfs.NewFake()
	const (
		pfDev  = "/sys/class/net/pf0/device"
		vfDev  = "/sys/class/net/vf0/device"
		pfBus  = "/sys/bus/pci/devices/0000:25:00.0"
		vf0Bus = "/sys/bus/pci/devices/0000:25:01.0"
		vf1Bus = "/sys/bus/pci/devices/0000:25:01.1"
	)

	// PF sysfs.
	fs.Files[pfDev+"/sriov_totalvfs"] = "64\n"
	fs.Files[pfDev+"/sriov_numvfs"] = "2\n"
	fs.Dirs[pfDev] = []string{"sriov_totalvfs", "sriov_numvfs", "virtfn0", "virtfn1", "driver"}
	fs.Links[pfDev+"/virtfn0"] = "../0000:25:01.0"
	fs.Links[pfDev+"/virtfn1"] = "../0000:25:01.1"

	// PF bus dir (used to resolve a VF's index back to the PF).
	fs.Dirs[pfBus] = []string{"virtfn0", "virtfn1"}
	fs.Links[pfBus+"/virtfn0"] = "../0000:25:01.0"
	fs.Links[pfBus+"/virtfn1"] = "../0000:25:01.1"

	// VF driver bindings.
	fs.Links[vf0Bus+"/driver"] = "../../../../bus/pci/drivers/iavf"
	fs.Links[vf1Bus+"/driver"] = "../../../../bus/pci/drivers/vfio-pci"

	// The VF0 netdev itself is a VF (physfn points back to the PF bus).
	fs.Links[vfDev+"/physfn"] = "../0000:25:00.0"
	fs.Dirs[vfDev] = []string{"physfn"}

	return fs
}

func TestSRIOVEnrichPF(t *testing.T) {
	fake := run.NewFake().SetResult(
		run.FakeResult{Stdout: string(testutil.Fixture(t, "ip/link-vf.json"))},
		"ip", "-details", "-json", "link", "show", "pf0",
	)
	ifaces := []model.Interface{
		{Name: "pf0", Bus: "0000:25:00.0", Driver: "ice", Type: model.TypePhysical},
		{Name: "vf0", Bus: "0000:25:01.0", Driver: "iavf", Type: model.TypePhysical},
	}

	NewSRIOV(fake, sriovFS()).Enrich(context.Background(), ifaces)

	pf := ifaces[0].SRIOV
	if pf == nil || !pf.Capable || pf.TotalVFs != 64 || pf.ConfiguredVFs != 2 || !pf.Enabled {
		t.Fatalf("pf sriov = %+v", pf)
	}
	if len(pf.VFs) != 2 {
		t.Fatalf("pf VFs = %d, want 2", len(pf.VFs))
	}

	vf0, vf1 := pf.VFs[0], pf.VFs[1]
	if vf0.Bus != "0000:25:01.0" || vf0.Driver != "iavf" || vf0.Netdev != "vf0" {
		t.Errorf("vf0 = %+v", vf0)
	}
	// Attributes come from the ip link fixture.
	if vf0.MAC != "aa:bb:cc:dd:ee:00" || vf0.VLAN != 100 || !vf0.SpoofCheck {
		t.Errorf("vf0 attrs = %+v", vf0)
	}
	if vf1.Driver != "vfio-pci" || vf1.Netdev != "" {
		t.Errorf("vf1 (vfio-pci, no netdev) = %+v", vf1)
	}
}

func TestSRIOVEnrichVF(t *testing.T) {
	ifaces := []model.Interface{
		{Name: "pf0", Bus: "0000:25:00.0", Type: model.TypePhysical},
		{Name: "vf0", Bus: "0000:25:01.0", Type: model.TypePhysical},
	}
	NewSRIOV(run.NewFake(), sriovFS()).Enrich(context.Background(), ifaces)

	vf := ifaces[1]
	if vf.Type != model.TypeVF {
		t.Errorf("vf0 type = %q, want vf", vf.Type)
	}
	if vf.SRIOV == nil || !vf.SRIOV.VF {
		t.Fatalf("vf0 sriov = %+v", vf.SRIOV)
	}
	if vf.SRIOV.PF != "pf0" || vf.SRIOV.PFBus != "0000:25:00.0" {
		t.Errorf("vf0 pf = %q (%q)", vf.SRIOV.PF, vf.SRIOV.PFBus)
	}
	if vf.SRIOV.VFIndex == nil || *vf.SRIOV.VFIndex != 0 {
		t.Errorf("vf0 index = %v, want 0", vf.SRIOV.VFIndex)
	}
}

func TestSRIOVEnrichNonSRIOV(t *testing.T) {
	ifaces := []model.Interface{{Name: "eth0", Bus: "0000:01:00.0"}}
	NewSRIOV(run.NewFake(), sysfs.NewFake()).Enrich(context.Background(), ifaces)
	if ifaces[0].SRIOV != nil {
		t.Errorf("non-SR-IOV interface should have nil SRIOV, got %+v", ifaces[0].SRIOV)
	}
}

func TestSRIOVTotalVFsZeroNotCapable(t *testing.T) {
	// sriov_totalvfs present but 0 means no VFs are possible; the device is not
	// SR-IOV capable and gets no SRIOV block.
	fs := sysfs.NewFake()
	fs.Files["/sys/class/net/eth0/device/sriov_totalvfs"] = "0\n"
	fs.Files["/sys/class/net/eth0/device/sriov_numvfs"] = "0\n"

	ifaces := []model.Interface{{Name: "eth0", Bus: "0000:01:00.0"}}
	NewSRIOV(run.NewFake(), fs).Enrich(context.Background(), ifaces)
	if ifaces[0].SRIOV != nil {
		t.Errorf("totalvfs=0 should yield no SRIOV block, got %+v", ifaces[0].SRIOV)
	}
}

func TestSRIOVPFNoConfiguredVFs(t *testing.T) {
	fs := sysfs.NewFake()
	fs.Files["/sys/class/net/pf0/device/sriov_totalvfs"] = "64\n"
	fs.Files["/sys/class/net/pf0/device/sriov_numvfs"] = "0\n"

	ifaces := []model.Interface{{Name: "pf0", Bus: "0000:25:00.0"}}
	NewSRIOV(run.NewFake(), fs).Enrich(context.Background(), ifaces)

	s := ifaces[0].SRIOV
	if s == nil || !s.Capable || s.Enabled || len(s.VFs) != 0 {
		t.Errorf("pf with 0 VFs = %+v", s)
	}
}
