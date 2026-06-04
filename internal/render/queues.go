package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// Queues renders the per-interface channel and ring table. Cells are
// "current/max"; a zero max renders blank (feature not applicable).
func (o Options) Queues(w io.Writer, queues []model.Queues) {
	headers := []string{"NAME", "COMBINED", "RX-CH", "TX-CH", "RX-RING", "TX-RING"}
	rows := make([][]string, 0, len(queues))
	for _, q := range queues {
		rows = append(rows, []string{
			q.Name,
			countCell(q.Combined),
			countCell(q.RxChannels),
			countCell(q.TxChannels),
			countCell(q.RxRing),
			countCell(q.TxRing),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// countCell renders a current/max pair, or blank when max is 0.
func countCell(c model.Count) string {
	if c.Max == 0 && c.Current == 0 {
		return ""
	}
	return strconv.Itoa(c.Current) + "/" + strconv.Itoa(c.Max)
}
