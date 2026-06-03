package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// Rules renders the routing policy rule table.
func (o Options) Rules(w io.Writer, rules []model.Rule) {
	headers := []string{"PRIORITY", "FROM", "TO", "IIF", "OIF", "FWMARK", "TABLE"}
	rows := make([][]string, 0, len(rules))
	for _, r := range rules {
		rows = append(rows, []string{
			strconv.Itoa(r.Priority),
			r.From,
			r.To,
			r.IIf,
			r.OIf,
			r.FWMark,
			r.Table,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
