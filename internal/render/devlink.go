package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// Devlink renders the devlink port table.
func (o Options) Devlink(w io.Writer, ports []model.DevlinkPort) {
	headers := []string{"HANDLE", "TYPE", "FLAVOUR", "NETDEV", "PF", "VF", "LANES"}
	rows := make([][]string, 0, len(ports))
	for _, p := range ports {
		lanes := ""
		if p.Lanes > 0 {
			lanes = strconv.Itoa(p.Lanes)
		}
		rows = append(rows, []string{
			p.Handle,
			p.Type,
			p.Flavour,
			p.Netdev,
			intPtrCell(p.PfNum),
			intPtrCell(p.VfNum),
			lanes,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// intPtrCell renders an optional int, blank when nil.
func intPtrCell(v *int) string {
	if v == nil {
		return ""
	}
	return strconv.Itoa(*v)
}
