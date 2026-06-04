package parse

import "testing"

func TestDevlinkPorts(t *testing.T) {
	data := []byte(`{"port":{
	  "pci/0000:25:00.0/0":{"type":"eth","netdev":"eth0","flavour":"physical","port":0,"lanes":4},
	  "pci/0000:25:00.0/1":{"type":"eth","netdev":"eth0v0","flavour":"pcivf","pfnum":0,"vfnum":1}
	}}`)
	ports, err := DevlinkPorts(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(ports) != 2 {
		t.Fatalf("ports = %d, want 2", len(ports))
	}
	// Sorted by handle: physical first.
	if ports[0].Flavour != "physical" || ports[0].Netdev != "eth0" || ports[0].Lanes != 4 {
		t.Errorf("port0 = %+v", ports[0])
	}
	vf := ports[1]
	if vf.Flavour != "pcivf" || vf.PfNum == nil || *vf.PfNum != 0 || vf.VfNum == nil || *vf.VfNum != 1 {
		t.Errorf("vf port = %+v", vf)
	}
}

func TestDevlinkPortsEmpty(t *testing.T) {
	ports, err := DevlinkPorts([]byte(`{"port":{}}`))
	if err != nil || len(ports) != 0 {
		t.Errorf("ports = %v, err = %v", ports, err)
	}
}
