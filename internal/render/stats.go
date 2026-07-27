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
	if hasRates(stats) {
		o.statsRates(w, stats)
		return
	}
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

func hasRates(stats []model.InterfaceStats) bool {
	for _, stat := range stats {
		if stat.Rates != nil || stat.RateStatus != "" {
			return true
		}
	}
	return false
}

func (o Options) statsRates(w io.Writer, stats []model.InterfaceStats) {
	headers := []string{"NAME", "RX RATE", "RX PPS", "RX ERR/s", "RX DROP/s", "TX RATE", "TX PPS", "TX ERR/s", "TX DROP/s", "UTIL"}
	rows := make([][]string, 0, len(stats))
	for _, stat := range stats {
		if stat.Rates == nil {
			rows = append(rows, []string{stat.Name, stat.RateStatus, "-", "-", "-", "-", "-", "-", "-", "-"})
			continue
		}
		rate := stat.Rates
		util := "-"
		if stat.LinkSpeedBps > 0 {
			util = fmt.Sprintf("RX %.1f%% / TX %.1f%%", rate.RxUtilization, rate.TxUtilization)
		}
		rows = append(rows, []string{
			stat.Name,
			humanBits(rate.RxBitsPerSec),
			humanNumber(rate.RxPacketsSec),
			humanNumber(rate.RxErrorsSec),
			humanNumber(rate.RxDroppedSec),
			humanBits(rate.TxBitsPerSec),
			humanNumber(rate.TxPacketsSec),
			humanNumber(rate.TxErrorsSec),
			humanNumber(rate.TxDroppedSec),
			util,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

func humanBits(rate float64) string {
	units := []string{"b/s", "Kb/s", "Mb/s", "Gb/s", "Tb/s"}
	unit := 0
	for rate >= 1000 && unit < len(units)-1 {
		rate /= 1000
		unit++
	}
	if rate >= 100 || unit == 0 {
		return fmt.Sprintf("%.0f %s", rate, units[unit])
	}
	return fmt.Sprintf("%.1f %s", rate, units[unit])
}

func humanNumber(rate float64) string {
	if rate == 0 {
		return "0"
	}
	if rate >= 100 {
		return fmt.Sprintf("%.0f", rate)
	}
	return fmt.Sprintf("%.1f", rate)
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
