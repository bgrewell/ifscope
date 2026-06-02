package parse

import (
	"encoding/json"

	"github.com/bgrewell/ifscope/internal/model"
)

// vfAttr mirrors one entry of a PF's vfinfo_list in `ip -details -json link`.
// Bools are pointers so an absent field is distinguishable from false.
type vfAttr struct {
	VF       int    `json:"vf"`
	Address  string `json:"address"`
	VlanList []struct {
		Vlan int `json:"vlan"`
	} `json:"vlan_list"`
	Vlan      *int   `json:"vlan"`
	Spoofchk  *bool  `json:"spoofchk"`
	LinkState string `json:"link_state"`
	Trust     *bool  `json:"trust"`
}

type ipLinkDetail struct {
	VFInfoList []vfAttr `json:"vfinfo_list"`
}

// VFAttrs parses `ip -details -json link show <pf>` and returns VF attributes
// keyed by VF index. Only the kernel-reported fields are populated; sysfs
// supplies bus/driver/netdev. A PF with no VFs yields an empty map.
func VFAttrs(data []byte) (map[int]model.VF, error) {
	var raw []ipLinkDetail
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	out := map[int]model.VF{}
	if len(raw) == 0 {
		return out, nil
	}
	for _, a := range raw[0].VFInfoList {
		vf := model.VF{
			Index:     a.VF,
			MAC:       a.Address,
			LinkState: a.LinkState,
		}
		switch {
		case len(a.VlanList) > 0:
			vf.VLAN = a.VlanList[0].Vlan
		case a.Vlan != nil:
			vf.VLAN = *a.Vlan
		}
		if a.Spoofchk != nil {
			vf.SpoofCheck = *a.Spoofchk
		}
		if a.Trust != nil {
			vf.Trust = *a.Trust
		}
		out[a.VF] = vf
	}
	return out, nil
}
