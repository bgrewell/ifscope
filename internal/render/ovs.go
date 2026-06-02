package render

import (
	"io"
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// OVS renders the Open vSwitch port topology, one row per port.
func (o Options) OVS(w io.Writer, ovs *model.OVS) {
	headers := []string{"BRIDGE", "PORT", "INTERFACES", "TAG", "TRUNKS", "VLANMODE"}
	var rows [][]string
	if ovs != nil {
		for _, p := range ovs.Ports {
			rows = append(rows, []string{
				p.Bridge,
				p.Name,
				strings.Join(p.Interfaces, "\n"),
				tagCell(p.Tag),
				intsCell(p.Trunks),
				p.VLANMode,
			})
		}
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// tagCell renders an optional access VLAN tag.
func tagCell(tag *int) string {
	if tag == nil {
		return ""
	}
	return strconv.Itoa(*tag)
}

// intsCell renders a list of ints comma-separated.
func intsCell(vals []int) string {
	if len(vals) == 0 {
		return ""
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = strconv.Itoa(v)
	}
	return strings.Join(parts, ",")
}
