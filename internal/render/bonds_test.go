package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
)

func TestBondsRender(t *testing.T) {
	bonds := []model.Bond{
		{
			Name:        "bond.les3",
			Mode:        "802.3ad (LACP)",
			State:       "UP",
			ActiveSlave: "enp67s0f2",
			Members: []model.BondMember{
				{Name: "enp67s0f2", State: "UP"},
				{Name: "enp67s0f3", State: "UP"},
			},
		},
		{
			Name:  "bond.les4",
			Mode:  "active-backup",
			State: "DOWN",
			Members: []model.BondMember{
				{Name: "enp67s0f0", State: "DOWN"},
				{Name: "enp67s0f1", State: "DOWN"},
			},
		},
	}

	var buf bytes.Buffer
	Options{Color: NewColor("never", false, nil)}.Bonds(&buf, bonds)
	out := buf.String()

	for _, want := range []string{
		"BOND", "MODE", "ACTIVE", "MEMBERS",
		"bond.les3", "802.3ad (LACP)", "enp67s0f2",
		"enp67s0f2 (up)", "enp67s0f3 (up)",
		"enp67s0f0 (down)",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q in:\n%s", want, out)
		}
	}

	// A bond with no active slave shows a dash.
	if !strings.Contains(out, "- ") {
		t.Errorf("expected '-' for empty active slave in:\n%s", out)
	}
}
