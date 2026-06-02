package correlate

import (
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
)

func TestAnnotateOVS(t *testing.T) {
	tag := 30
	ovs := &model.OVS{
		Ports: []model.OVSPort{
			{Name: "vlan30", Bridge: "br-int", Interfaces: []string{"vlan30"}, Tag: &tag},
			{Name: "eth0", Bridge: "br-int", Interfaces: []string{"eth0"}},
		},
	}
	ifaces := []model.Interface{
		{Name: "eth0"},
		{Name: "vlan30"},
		{Name: "eth1"}, // not in OVS
	}

	AnnotateOVS(ovs, ifaces)

	if ifaces[0].OVS == nil || ifaces[0].OVS.Bridge != "br-int" {
		t.Errorf("eth0 membership = %+v", ifaces[0].OVS)
	}
	if ifaces[1].OVS == nil || ifaces[1].OVS.Tag == nil || *ifaces[1].OVS.Tag != 30 {
		t.Errorf("vlan30 membership = %+v", ifaces[1].OVS)
	}
	if ifaces[2].OVS != nil {
		t.Errorf("eth1 should not be annotated, got %+v", ifaces[2].OVS)
	}
}

func TestAnnotateOVSNil(t *testing.T) {
	ifaces := []model.Interface{{Name: "eth0"}}
	AnnotateOVS(nil, ifaces) // must not panic
	if ifaces[0].OVS != nil {
		t.Errorf("nil OVS should leave interfaces untouched")
	}
}
