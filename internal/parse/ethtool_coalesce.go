package parse

import (
	"regexp"
	"strconv"
	"strings"
)

// Coalesce holds the interrupt-coalescing settings ifscope surfaces.
type Coalesce struct {
	RxUsecs    int
	TxUsecs    int
	AdaptiveRx bool
	AdaptiveTx bool
}

var adaptiveRe = regexp.MustCompile(`Adaptive RX:\s*(\w+)\s+TX:\s*(\w+)`)

// EthtoolCoalesce parses `ethtool -c` interrupt-coalescing output.
func EthtoolCoalesce(data []byte) Coalesce {
	var c Coalesce
	text := string(data)
	if m := adaptiveRe.FindStringSubmatch(text); m != nil {
		c.AdaptiveRx = m[1] == "on"
		c.AdaptiveTx = m[2] == "on"
	}
	for _, raw := range strings.Split(text, "\n") {
		key, val, ok := strings.Cut(strings.TrimSpace(raw), ":")
		if !ok {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(val))
		if err != nil {
			continue
		}
		switch strings.TrimSpace(key) {
		case "rx-usecs":
			c.RxUsecs = n
		case "tx-usecs":
			c.TxUsecs = n
		}
	}
	return c
}

var rssRingsRe = regexp.MustCompile(`with (\d+) RX ring`)

// EthtoolRSSRings parses the ring count from `ethtool -x`'s header
// ("... with N RX ring(s)"); 0 if absent.
func EthtoolRSSRings(data []byte) int {
	if m := rssRingsRe.FindSubmatch(data); m != nil {
		if n, err := strconv.Atoi(string(m[1])); err == nil {
			return n
		}
	}
	return 0
}
