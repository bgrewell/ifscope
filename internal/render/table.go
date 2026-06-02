package render

import (
	"io"
	"strings"
	"unicode/utf8"
)

// Table is a renderable grid. A cell may contain newlines to span multiple
// display lines (used for address and altname lists).
type Table struct {
	Headers []string
	Rows    [][]string
}

// WriteUnicode renders the table with box-drawing borders.
func (t Table) WriteUnicode(w io.Writer) {
	widths := t.columnWidths()

	writeBorder(w, widths, "┌", "┬", "┐")
	writeRow(w, t.Headers, widths)
	writeBorder(w, widths, "├", "┼", "┤")
	for _, row := range t.Rows {
		writeRow(w, row, widths)
	}
	writeBorder(w, widths, "└", "┴", "┘")
}

// WriteBarebones renders the table as plain " | "-delimited rows with no
// borders, suitable for copy/paste and downstream text processing.
func (t Table) WriteBarebones(w io.Writer) {
	widths := t.columnWidths()
	writeBareRow(w, t.Headers, widths)
	for _, row := range t.Rows {
		writeBareRow(w, row, widths)
	}
}

// columnWidths returns the max visible width per column across header and cells.
func (t Table) columnWidths() []int {
	widths := make([]int, len(t.Headers))
	for i, h := range t.Headers {
		widths[i] = visibleWidth(h)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if i >= len(widths) {
				continue
			}
			for _, line := range strings.Split(cell, "\n") {
				if wpx := visibleWidth(line); wpx > widths[i] {
					widths[i] = wpx
				}
			}
		}
	}
	return widths
}

func writeBorder(w io.Writer, widths []int, left, mid, right string) {
	var b strings.Builder
	b.WriteString(left)
	for i, width := range widths {
		b.WriteString(strings.Repeat("─", width+2))
		if i < len(widths)-1 {
			b.WriteString(mid)
		}
	}
	b.WriteString(right)
	b.WriteByte('\n')
	io.WriteString(w, b.String())
}

// writeRow renders a single logical row, expanding multi-line cells to as many
// physical lines as the tallest cell requires.
func writeRow(w io.Writer, cells []string, widths []int) {
	lines := splitCells(cells, len(widths))
	height := 1
	for _, cl := range lines {
		if len(cl) > height {
			height = len(cl)
		}
	}
	for line := 0; line < height; line++ {
		var b strings.Builder
		b.WriteString("│")
		for col, width := range widths {
			text := ""
			if col < len(lines) && line < len(lines[col]) {
				text = lines[col][line]
			}
			b.WriteString(" ")
			b.WriteString(pad(text, width))
			b.WriteString(" │")
		}
		b.WriteByte('\n')
		io.WriteString(w, b.String())
	}
}

func writeBareRow(w io.Writer, cells []string, widths []int) {
	lines := splitCells(cells, len(widths))
	height := 1
	for _, cl := range lines {
		if len(cl) > height {
			height = len(cl)
		}
	}
	for line := 0; line < height; line++ {
		parts := make([]string, len(widths))
		for col, width := range widths {
			text := ""
			if col < len(lines) && line < len(lines[col]) {
				text = lines[col][line]
			}
			parts[col] = pad(text, width)
		}
		io.WriteString(w, strings.TrimRight(strings.Join(parts, " | "), " ")+"\n")
	}
}

// splitCells splits each cell into its display lines, padding the slice to n
// columns.
func splitCells(cells []string, n int) [][]string {
	out := make([][]string, n)
	for i := 0; i < n; i++ {
		if i < len(cells) {
			out[i] = strings.Split(cells[i], "\n")
		} else {
			out[i] = []string{""}
		}
	}
	return out
}

// pad right-pads s with spaces to the given visible width.
func pad(s string, width int) string {
	gap := width - visibleWidth(s)
	if gap <= 0 {
		return s
	}
	return s + strings.Repeat(" ", gap)
}

// visibleWidth returns the rune count of s ignoring ANSI SGR escape sequences,
// so colored cells align correctly.
func visibleWidth(s string) int {
	count := 0
	for i := 0; i < len(s); {
		if s[i] == '\033' {
			// Skip an escape sequence: ESC [ ... letter.
			j := i + 1
			if j < len(s) && s[j] == '[' {
				j++
				for j < len(s) && !isAnsiFinal(s[j]) {
					j++
				}
				if j < len(s) {
					j++ // consume the final byte
				}
				i = j
				continue
			}
		}
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
		count++
	}
	return count
}

func isAnsiFinal(b byte) bool {
	return b >= '@' && b <= '~'
}
