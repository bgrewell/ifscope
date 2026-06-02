package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// SRIOV renders the SR-IOV table: one row per VF, or a single row per PF that
// has no configured VFs. Interfaces without SR-IOV are omitted.
func (o Options) SRIOV(w io.Writer, ifaces []model.Interface) {
	headers := []string{"PF", "PF BUS", "DRIVER", "TOTAL", "CFG", "VF", "VF BUS", "VF DRIVER", "VF NETDEV", "MAC", "VLAN", "SPOOF", "TRUST", "LINK"}
	var rows [][]string

	for _, i := range ifaces {
		s := i.SRIOV
		if s == nil || s.VF || !s.Capable {
			continue
		}
		total := strconv.Itoa(s.TotalVFs)
		cfg := strconv.Itoa(s.ConfiguredVFs)
		if len(s.VFs) == 0 {
			rows = append(rows, []string{i.Name, busDisplay(i.Bus), i.Driver, total, cfg, "", "", "", "", "", "", "", "", ""})
			continue
		}
		for _, vf := range s.VFs {
			rows = append(rows, []string{
				i.Name, busDisplay(i.Bus), i.Driver, total, cfg,
				strconv.Itoa(vf.Index),
				busDisplay(vf.Bus),
				vf.Driver,
				vf.Netdev,
				vf.MAC,
				vlanCell(vf.VLAN),
				onOff(vf.SpoofCheck),
				onOff(vf.Trust),
				vf.LinkState,
			})
		}
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

func vlanCell(v int) string {
	if v == 0 {
		return ""
	}
	return strconv.Itoa(v)
}

func onOff(b bool) string {
	if b {
		return "on"
	}
	return "off"
}
