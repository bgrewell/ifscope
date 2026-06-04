package parse

import "testing"

func TestEthtoolFeatures(t *testing.T) {
	data := []byte(`Features for eth0:
rx-checksumming: on
tx-checksumming: on
	tx-checksum-ipv4: on
	tx-checksum-ip-generic: off [fixed]
scatter-gather: on
tcp-segmentation-offload: on
generic-receive-offload: on
large-receive-offload: off [fixed]`)
	f := EthtoolFeatures(data)
	if f["rx-checksumming"] != "on" {
		t.Errorf("rx-checksumming = %q", f["rx-checksumming"])
	}
	if f["large-receive-offload"] != "off [fixed]" {
		t.Errorf("lro = %q", f["large-receive-offload"])
	}
	// Sub-features (indented) and the header must be excluded.
	if _, ok := f["tx-checksum-ipv4"]; ok {
		t.Error("sub-feature should be excluded")
	}
	if _, ok := f["Features for eth0"]; ok {
		t.Error("header should be excluded")
	}
}
