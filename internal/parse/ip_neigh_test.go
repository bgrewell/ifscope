package parse

import "testing"

func TestIPNeighbors(t *testing.T) {
	data := []byte(`[
	  {"dst":"172.26.1.1","dev":"eth0","lladdr":"aa:bb:cc:dd:ee:ff","state":["REACHABLE"]},
	  {"dst":"fe80::1","dev":"eth0","lladdr":"aa:bb:cc:dd:ee:00","state":["STALE"],"router":true},
	  {"dst":"172.26.1.9","dev":"eth0","state":["FAILED"]}
	]`)
	ns, err := IPNeighbors(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(ns) != 3 {
		t.Fatalf("neighbors = %d, want 3", len(ns))
	}
	if ns[0].State != "REACHABLE" || ns[0].Family != "inet" || ns[0].LLAddr != "aa:bb:cc:dd:ee:ff" {
		t.Errorf("n0 = %+v", ns[0])
	}
	if ns[1].Family != "inet6" || !ns[1].Router {
		t.Errorf("n1 = %+v", ns[1])
	}
	if ns[2].LLAddr != "" || ns[2].State != "FAILED" {
		t.Errorf("incomplete neighbor = %+v", ns[2])
	}
}
