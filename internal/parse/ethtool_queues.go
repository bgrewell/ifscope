package parse

import (
	"strconv"
	"strings"

	"github.com/bgrewell/ifscope/internal/model"
)

// ethtoolMaxCur parses ethtool's two-section "Pre-set maximums" / "Current
// hardware settings" layout (shared by -l channels and -g rings) into
// label→value maps. Non-numeric values ("n/a") are skipped.
func ethtoolMaxCur(data []byte) (max, cur map[string]int) {
	max, cur = map[string]int{}, map[string]int{}
	target := max
	for _, raw := range strings.Split(string(data), "\n") {
		t := strings.TrimSpace(raw)
		switch {
		case strings.HasPrefix(t, "Pre-set maximums"):
			target = max
			continue
		case strings.HasPrefix(t, "Current hardware settings"):
			target = cur
			continue
		}
		key, val, ok := strings.Cut(t, ":")
		if !ok {
			continue
		}
		if n, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			target[strings.TrimSpace(key)] = n
		}
	}
	return max, cur
}

// EthtoolChannels parses `ethtool -l` into combined/rx/tx channel counts.
func EthtoolChannels(data []byte) (combined, rx, tx model.Count) {
	max, cur := ethtoolMaxCur(data)
	combined = model.Count{Current: cur["Combined"], Max: max["Combined"]}
	rx = model.Count{Current: cur["RX"], Max: max["RX"]}
	tx = model.Count{Current: cur["TX"], Max: max["TX"]}
	return combined, rx, tx
}

// EthtoolRings parses `ethtool -g` into rx/tx ring sizes.
func EthtoolRings(data []byte) (rx, tx model.Count) {
	max, cur := ethtoolMaxCur(data)
	rx = model.Count{Current: cur["RX"], Max: max["RX"]}
	tx = model.Count{Current: cur["TX"], Max: max["TX"]}
	return rx, tx
}
