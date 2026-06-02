package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
)

func TestVisibleWidthIgnoresANSI(t *testing.T) {
	colored := ansiGreen + "UP" + ansiReset
	if got := visibleWidth(colored); got != 2 {
		t.Errorf("visibleWidth(colored UP) = %d, want 2", got)
	}
	if got := visibleWidth("hello"); got != 5 {
		t.Errorf("visibleWidth(hello) = %d, want 5", got)
	}
}

func TestColorMode(t *testing.T) {
	if NewColor("never", false, nil).State("UP") != "UP" {
		t.Error("never mode should not color")
	}
	got := NewColor("always", false, nil).State("UP")
	if !strings.Contains(got, ansiGreen) {
		t.Error("always mode should color UP green")
	}
	if NewColor("always", true, nil).State("UP") != "UP" {
		t.Error("no-color override should win over always")
	}
}

func TestBarebonesTable(t *testing.T) {
	tbl := Table{
		Headers: []string{"NAME", "STATE"},
		Rows:    [][]string{{"eth0", "UP"}, {"eth1", "DOWN"}},
	}
	var buf bytes.Buffer
	tbl.WriteBarebones(&buf)
	out := buf.String()

	if !strings.Contains(out, "NAME | STATE") {
		t.Errorf("missing header row in:\n%s", out)
	}
	if !strings.Contains(out, "eth0 | UP") {
		t.Errorf("missing data row in:\n%s", out)
	}
}

func TestUnicodeTableHasBorders(t *testing.T) {
	tbl := Table{Headers: []string{"A"}, Rows: [][]string{{"x"}}}
	var buf bytes.Buffer
	tbl.WriteUnicode(&buf)
	out := buf.String()
	for _, want := range []string{"┌", "┐", "│", "└", "┘"} {
		if !strings.Contains(out, want) {
			t.Errorf("unicode table missing %q in:\n%s", want, out)
		}
	}
}

func TestInterfacesTableMultilineAddrs(t *testing.T) {
	ifaces := []model.Interface{{
		ID:    2,
		Name:  "eth0",
		State: "UP",
		Addresses: []model.Address{
			{Family: "inet", Local: "192.0.2.10", PrefixLen: 24},
			{Family: "inet", Local: "192.0.2.11", PrefixLen: 24},
			{Family: "inet6", Local: "fe80::1", PrefixLen: 64},
		},
	}}

	var buf bytes.Buffer
	Options{Color: NewColor("never", false, nil)}.Interfaces(&buf, ifaces)
	out := buf.String()

	if !strings.Contains(out, "192.0.2.10/24") || !strings.Contains(out, "192.0.2.11/24") {
		t.Errorf("expected both IPv4 addresses in:\n%s", out)
	}
	if strings.Contains(out, "fe80::1") {
		t.Errorf("IPv6 should be excluded from default table:\n%s", out)
	}
}
