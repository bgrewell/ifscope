package render

import (
	"io"

	"github.com/bgrewell/ifscope/internal/model"
)

// Tests renders the connectivity test results.
func (o Options) Tests(w io.Writer, tests []model.TestResult) {
	headers := []string{"TEST", "STATUS", "TARGET", "LATENCY", "DETAILS"}
	rows := make([][]string, 0, len(tests))
	for _, t := range tests {
		details := t.Details
		if t.Error != "" {
			details = t.Error
		}
		rows = append(rows, []string{
			t.Name,
			o.Color.TestStatus(string(t.Status)),
			t.Target,
			t.Latency,
			details,
		})
	}
	o.write(w, Table{Headers: headers, Rows: rows})
}
