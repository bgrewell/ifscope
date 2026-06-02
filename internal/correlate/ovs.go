package correlate

import "github.com/bgrewell/ifscope/internal/model"

// AnnotateOVS tags interfaces and VLANs with their OVS bridge/port membership,
// matching by interface name. Interfaces not present in OVS are left untouched.
func AnnotateOVS(ovs *model.OVS, groups ...[]model.Interface) {
	if ovs == nil {
		return
	}

	membership := map[string]model.OVSMembership{}
	for _, p := range ovs.Ports {
		for _, ifName := range p.Interfaces {
			membership[ifName] = model.OVSMembership{Bridge: p.Bridge, Port: p.Name, Tag: p.Tag}
		}
		// The port itself (e.g. an internal port sharing the port name) maps too.
		if _, ok := membership[p.Name]; !ok {
			membership[p.Name] = model.OVSMembership{Bridge: p.Bridge, Port: p.Name, Tag: p.Tag}
		}
	}

	for _, group := range groups {
		for i := range group {
			if m, ok := membership[group[i].Name]; ok {
				mc := m
				group[i].OVS = &mc
			}
		}
	}
}
