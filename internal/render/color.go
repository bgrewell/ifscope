// Package render turns a model.Report into human tables or JSON. Output
// concerns (color, borders, column widths) are isolated here so collection and
// correlation logic stay presentation-free.
package render

import (
	"os"
)

// ANSI SGR sequences used for state coloring.
const (
	ansiReset = "\033[0m"
	ansiRed   = "\033[31m"
	ansiGreen = "\033[32m"
)

// Color applies optional ANSI coloring. The zero value is disabled.
type Color struct {
	enabled bool
}

// NewColor resolves whether coloring should be active from the color mode
// ("auto"|"always"|"never"), an explicit no-color override, and whether w is a
// terminal (consulted only in "auto" mode).
func NewColor(mode string, noColor bool, w *os.File) Color {
	if noColor || mode == "never" {
		return Color{enabled: false}
	}
	if mode == "always" {
		return Color{enabled: true}
	}
	return Color{enabled: isTerminal(w)}
}

// State colors an interface state: green for UP, red for down variants.
func (c Color) State(state string) string {
	if !c.enabled {
		return state
	}
	switch state {
	case "UP":
		return ansiGreen + state + ansiReset
	case "DOWN", "LOWERLAYERDOWN":
		return ansiRed + state + ansiReset
	default:
		return state
	}
}

// isTerminal reports whether w refers to a character device (a TTY).
func isTerminal(w *os.File) bool {
	if w == nil {
		return false
	}
	fi, err := w.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
