package parse

import "testing"

func TestEthtoolChannels(t *testing.T) {
	data := []byte(`Channel parameters for eth0:
Pre-set maximums:
RX:		192
TX:		192
Other:		1
Combined:	192
Current hardware settings:
RX:		0
TX:		0
Other:		1
Combined:	48`)
	combined, rx, tx := EthtoolChannels(data)
	if combined.Current != 48 || combined.Max != 192 {
		t.Errorf("combined = %+v, want 48/192", combined)
	}
	if rx.Current != 0 || tx.Max != 192 {
		t.Errorf("rx=%+v tx=%+v", rx, tx)
	}
}

func TestEthtoolRings(t *testing.T) {
	data := []byte(`Ring parameters for eth0:
Pre-set maximums:
RX:			8160
RX Mini:		n/a
TX:			8160
Current hardware settings:
RX:			2048
RX Mini:		n/a
TX:			256`)
	rx, tx := EthtoolRings(data)
	if rx.Current != 2048 || rx.Max != 8160 {
		t.Errorf("rx ring = %+v, want 2048/8160", rx)
	}
	if tx.Current != 256 || tx.Max != 8160 {
		t.Errorf("tx ring = %+v, want 256/8160", tx)
	}
}
