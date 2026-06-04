package parse

import (
	"strconv"
	"strings"
)

// IRQLine is a parsed /proc/interrupts row: the IRQ number and its name (the
// trailing label, e.g. "ice-enp23s0np0-TxRx-0").
type IRQLine struct {
	Number int
	Name   string
}

// Interrupts parses /proc/interrupts. Non-numeric IRQ rows (NMI, LOC, ...) are
// skipped. The per-CPU counts between the number and the trailing name are
// ignored; the name is the last whitespace field.
func Interrupts(data []byte) []IRQLine {
	var out []IRQLine
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		head := strings.TrimSuffix(fields[0], ":")
		if head == fields[0] { // no trailing colon
			continue
		}
		num, err := strconv.Atoi(head)
		if err != nil {
			continue
		}
		out = append(out, IRQLine{Number: num, Name: fields[len(fields)-1]})
	}
	return out
}
