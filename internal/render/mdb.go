package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// MDB renders the bridge multicast-database table.
func (o Options) MDB(w io.Writer, entries []model.MDBEntry) {
	headers := []string{"BRIDGE", "PORT", "GROUP", "VLAN", "STATE"}
	rows := make([][]string, 0, len(entries))
	for _, e := range entries {
		vlan := ""
		if e.VLAN > 0 {
			vlan = strconv.Itoa(e.VLAN)
		}
		rows = append(rows, []string{e.Bridge, e.Port, e.Group, vlan, e.State})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
