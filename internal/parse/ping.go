package parse

import "regexp"

// rttRe captures the average from ping's "rtt min/avg/max/mdev = a/b/c/d ms".
var rttRe = regexp.MustCompile(`rtt[^=]*=\s*[0-9.]+/([0-9.]+)/`)

// PingAvgMillis returns ping's average round-trip time as a string of
// milliseconds (e.g. "6.154"), and whether it was found. Absence of a summary
// line indicates the ping produced no replies.
func PingAvgMillis(data []byte) (string, bool) {
	if m := rttRe.FindSubmatch(data); m != nil {
		return string(m[1]), true
	}
	return "", false
}
