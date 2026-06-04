package render

import (
	"io"

	"github.com/bgrewell/ifscope/internal/model"
)

// Multicast renders the IP multicast group-membership table.
func (o Options) Multicast(w io.Writer, groups []model.MulticastGroup) {
	headers := []string{"INTERFACE", "FAMILY", "ADDRESS"}
	rows := make([][]string, 0, len(groups))
	for _, g := range groups {
		rows = append(rows, []string{g.Interface, g.Family, g.Address})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
