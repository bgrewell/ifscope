package render

import (
	"fmt"
	"io"
	"strconv"

	"github.com/bgrewell/ifscope/internal/model"
)

// Stats renders the per-interface counter table. Byte counts are humanized;
// the exact values remain in JSON output.
func (o Options) Stats(w io.Writer, stats []model.InterfaceStats) {
	headers := []string{"NAME", "RX BYTES", "RX PKTS", "RX ERR", "RX DROP", "TX BYTES", "TX PKTS", "TX ERR", "TX DROP"}
	rows := make([][]string, 0, len(stats))
	for _, s := range stats {
		rows = append(rows, []string{
			s.Name,
			humanBytes(s.RxBytes),
			strconv.FormatUint(s.RxPackets, 10),
			strconv.FormatUint(s.RxErrors, 10),
			strconv.FormatUint(s.RxDropped, 10),
			humanBytes(s.TxBytes),
			strconv.FormatUint(s.TxPackets, 10),
			strconv.FormatUint(s.TxErrors, 10),
			strconv.FormatUint(s.TxDropped, 10),
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// humanBytes formats a byte count with binary (1024) units.
func humanBytes(n uint64) string {
	const unit = 1024
	if n < unit {
		return strconv.FormatUint(n, 10) + " B"
	}
	div, exp := uint64(unit), 0
	for n/div >= unit && exp < 5 {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}
