package parse

import "testing"

func TestEthtoolTimestamp(t *testing.T) {
	data := []byte(`Time stamping parameters for eth0:
Capabilities:
	hardware-transmit
	software-transmit
	hardware-receive
	software-receive
	software-system-clock
Hardware timestamp provider index: 4
Hardware Transmit Timestamp Modes:
	off
	on`)
	p := EthtoolTimestamp(data)
	if !p.HWTx || !p.HWRx || !p.SWTx || !p.SWRx {
		t.Errorf("caps = %+v", p)
	}
	if p.PHCIndex == nil || *p.PHCIndex != 4 {
		t.Errorf("phc = %v, want 4", p.PHCIndex)
	}
}

func TestEthtoolTimestampSoftwareOnly(t *testing.T) {
	data := []byte(`Time stamping parameters for veth0:
Capabilities:
	software-transmit
	software-receive
	software-system-clock
PTP Hardware Clock: none`)
	p := EthtoolTimestamp(data)
	if p.HWTx || p.HWRx {
		t.Errorf("should have no hw caps: %+v", p)
	}
	if p.PHCIndex != nil {
		t.Errorf("phc should be nil (none), got %v", p.PHCIndex)
	}
	if !p.SWTx {
		t.Errorf("sw-tx should be set")
	}
}
