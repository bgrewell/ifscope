package render

import (
	"io"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// offloadColumns maps display headers to ethtool -k feature names.
var offloadColumns = []struct{ header, feature string }{
	{"RX-CSUM", "rx-checksumming"},
	{"TX-CSUM", "tx-checksumming"},
	{"SG", "scatter-gather"},
	{"TSO", "tcp-segmentation-offload"},
	{"GSO", "generic-segmentation-offload"},
	{"GRO", "generic-receive-offload"},
	{"LRO", "large-receive-offload"},
}

// Offloads renders a per-interface summary of common NIC offloads. Full feature
// sets are available in JSON output.
func (o Options) Offloads(w io.Writer, offloads []model.Offloads) {
	headers := []string{"NAME"}
	for _, c := range offloadColumns {
		headers = append(headers, c.header)
	}
	rows := make([][]string, 0, len(offloads))
	for _, off := range offloads {
		row := []string{off.Name}
		for _, c := range offloadColumns {
			row = append(row, offloadState(off.Features[c.feature]))
		}
		rows = append(rows, row)
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}

// offloadState reduces a feature value to on/off, marking fixed with "*".
func offloadState(v string) string {
	if v == "" {
		return ""
	}
	fixed := strings.Contains(v, "[fixed]")
	state := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(v), "[fixed]"))
	if fixed {
		return state + "*"
	}
	return state
}
