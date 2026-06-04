package render

import (
	"io"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// LLDP renders the LLDP neighbor table.
func (o Options) LLDP(w io.Writer, neighbors []model.LLDPNeighbor) {
	headers := []string{"LOCAL", "CHASSIS", "PORT", "PORT DESCR", "MGMT", "CAPS", "TTL"}
	rows := make([][]string, 0, len(neighbors))
	for _, n := range neighbors {
		rows = append(rows, []string{
			n.LocalPort,
			n.Chassis,
			n.PortID,
			n.PortDescr,
			strings.Join(n.MgmtIPs, "\n"),
			strings.Join(n.Capabilities, ","),
			n.TTL,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
