package parse

import "testing"

func TestEthtoolCoalesce(t *testing.T) {
	data := []byte(`Coalesce parameters for eth0:
Adaptive RX: on  TX: off
rx-usecs:	50
tx-usecs:	25
rx-frames:	0`)
	c := EthtoolCoalesce(data)
	if !c.AdaptiveRx || c.AdaptiveTx {
		t.Errorf("adaptive rx=%v tx=%v", c.AdaptiveRx, c.AdaptiveTx)
	}
	if c.RxUsecs != 50 || c.TxUsecs != 25 {
		t.Errorf("usecs rx=%d tx=%d", c.RxUsecs, c.TxUsecs)
	}
}

func TestEthtoolRSSRings(t *testing.T) {
	data := []byte("RX flow hash indirection table for eth0 with 48 RX ring(s):\n    0:  0  1  2  3\n")
	if got := EthtoolRSSRings(data); got != 48 {
		t.Errorf("rss rings = %d, want 48", got)
	}
	if got := EthtoolRSSRings([]byte("no match")); got != 0 {
		t.Errorf("no-match = %d, want 0", got)
	}
}
