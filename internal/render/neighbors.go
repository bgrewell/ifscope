package render

import (
	"io"

	"github.com/bgrewell/ifscope/internal/model"
)

// Neighbors renders the ARP/NDP neighbor table.
func (o Options) Neighbors(w io.Writer, neighbors []model.Neighbor) {
	headers := []string{"DST", "DEV", "LLADDR", "STATE"}
	rows := make([][]string, 0, len(neighbors))
	for _, n := range neighbors {
		rows = append(rows, []string{n.Dst, n.Dev, n.LLAddr, n.State})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
