package render

import (
	"fmt"
	"io"

	"github.com/bgrewell/ifscope/internal/model"
)

// TCClass renders the traffic-control shaping-class table.
func (o Options) TCClass(w io.Writer, classes []model.TCClass) {
	headers := []string{"DEV", "KIND", "HANDLE", "PARENT", "RATE", "CEIL"}
	rows := make([][]string, 0, len(classes))
	for _, c := range classes {
		parent := c.Parent
		if c.Root {
			parent = "root"
		}
		rows = append(rows, []string{
			c.Dev, c.Kind, c.Handle, parent, humanBitrate(c.Rate), humanBitrate(c.Ceil),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// humanBitrate formats bits/s with decimal (1000) units.
func humanBitrate(bps uint64) string {
	if bps == 0 {
		return ""
	}
	const unit = 1000
	if bps < unit {
		return fmt.Sprintf("%d bit", bps)
	}
	div, exp := uint64(unit), 0
	for bps/div >= unit && exp < 4 {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.3g %cbit", float64(bps)/float64(div), "kMGT"[exp])
}
