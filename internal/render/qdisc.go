package render

import (
	"io"

	"github.com/bgrewell/ifscope/internal/model"
)

// Qdisc renders the per-device root queueing-discipline table.
func (o Options) Qdisc(w io.Writer, qdiscs []model.Qdisc) {
	headers := []string{"DEV", "KIND", "HANDLE"}
	rows := make([][]string, 0, len(qdiscs))
	for _, q := range qdiscs {
		rows = append(rows, []string{q.Dev, q.Kind, q.Handle})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
