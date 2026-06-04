package render

import (
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// Queues renders the per-interface channel/ring/coalesce/steering table. Count
// cells are "current/max" (blank when not applicable).
func (o Options) Queues(w io.Writer, queues []model.Queues) {
	headers := []string{"NAME", "COMBINED", "RX-RING", "TX-RING", "COALESCE", "RSS", "RPS", "XPS"}
	rows := make([][]string, 0, len(queues))
	for _, q := range queues {
		rows = append(rows, []string{
			q.Name,
			countCell(q.Combined),
			countCell(q.RxRing),
			countCell(q.TxRing),
			coalesceCell(q),
			intCell(q.RSSRings),
			intCell(q.RPSQueues),
			intCell(q.XPSQueues),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// coalesceCell renders "rxusec/txusec", with " (a)" appended when adaptive.
func coalesceCell(q model.Queues) string {
	if q.RxUsecs == 0 && q.TxUsecs == 0 && !q.AdaptiveRx && !q.AdaptiveTx {
		return ""
	}
	s := strconv.Itoa(q.RxUsecs) + "/" + strconv.Itoa(q.TxUsecs)
	if q.AdaptiveRx || q.AdaptiveTx {
		s += " (a)"
	}
	return s
}

// intCell renders a non-zero int, blank for zero.
func intCell(n int) string {
	if n == 0 {
		return ""
	}
	return strconv.Itoa(n)
}

// countCell renders a current/max pair, or blank when max is 0.
func countCell(c model.Count) string {
	if c.Max == 0 && c.Current == 0 {
		return ""
	}
	return strconv.Itoa(c.Current) + "/" + strconv.Itoa(c.Max)
}
