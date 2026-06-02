package render

import (
	"encoding/json"
	"fmt"
	"io"
)

// JSON writes v as JSON to w. When pretty is set, output is indented. JSON is
// the only thing written to w in JSON mode; warnings travel inside the value
// (Report.Warnings), never mixed into the stream.
func JSON(w io.Writer, v any, pretty bool) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	}
	if err := enc.Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
