package render

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bgrewell/ifscope/internal/model"
)

func TestHumanBytes(t *testing.T) {
	cases := map[uint64]string{
		0:          "0 B",
		512:        "512 B",
		1024:       "1.0 KiB",
		1536:       "1.5 KiB",
		1048576:    "1.0 MiB",
		1073741824: "1.0 GiB",
	}
	for in, want := range cases {
		if got := humanBytes(in); got != want {
			t.Errorf("humanBytes(%d) = %q, want %q", in, got, want)
		}
	}
}

func TestStatsRendersRates(t *testing.T) {
	var out bytes.Buffer
	Options{Barebones: true}.Stats(&out, []model.InterfaceStats{{
		Name:         "eth0",
		LinkSpeedBps: 1_000_000_000,
		Rates: &model.InterfaceStatsRate{
			RxBitsPerSec:  12_500_000,
			TxBitsPerSec:  2_000,
			RxUtilization: 1.25,
			TxUtilization: 0.0002,
		},
	}})
	got := out.String()
	for _, want := range []string{"RX RATE", "12.5 Mb/s", "2.0 Kb/s", "RX 1.2% / TX 0.0%"} {
		if !strings.Contains(got, want) {
			t.Errorf("output does not contain %q:\n%s", want, got)
		}
	}
}
