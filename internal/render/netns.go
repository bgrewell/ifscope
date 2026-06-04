package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// Netns renders the network namespace table.
func (o Options) Netns(w io.Writer, namespaces []model.Netns) {
	headers := []string{"NAME", "ID", "INTERFACES"}
	rows := make([][]string, 0, len(namespaces))
	for _, n := range namespaces {
		id := ""
		if n.ID != nil {
			id = strconv.Itoa(*n.ID)
		}
		ifaces := ""
		if n.Interfaces != nil {
			ifaces = strconv.Itoa(*n.Interfaces)
		}
		rows = append(rows, []string{n.Name, id, ifaces})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
