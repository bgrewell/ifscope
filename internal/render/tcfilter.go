package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// TCFilter renders the traffic-control filter (classifier) table.
func (o Options) TCFilter(w io.Writer, filters []model.TCFilter) {
	headers := []string{"DEV", "PARENT", "PRIO", "PROTO", "KIND", "FLOWID"}
	rows := make([][]string, 0, len(filters))
	for _, f := range filters {
		pref := ""
		if f.Pref > 0 {
			pref = strconv.Itoa(f.Pref)
		}
		rows = append(rows, []string{f.Dev, f.Parent, pref, f.Protocol, f.Kind, f.FlowID})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
