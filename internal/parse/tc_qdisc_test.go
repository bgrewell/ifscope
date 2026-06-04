package parse

import "testing"

func TestQdiscs(t *testing.T) {
	// Root mq on eth0 with per-queue pfifo_fast children, and a shaped eth1.
	data := []byte(`[
	  {"kind":"noqueue","handle":"0:","dev":"lo","root":true},
	  {"kind":"mq","handle":"0:","dev":"eth0","root":true},
	  {"kind":"pfifo_fast","handle":"0:","dev":"eth0","parent":":1"},
	  {"kind":"fq_codel","handle":"8001:","dev":"eth1","root":true}
	]`)
	qs, err := Qdiscs(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(qs) != 4 {
		t.Fatalf("qdiscs = %d, want 4", len(qs))
	}
	roots := 0
	for _, q := range qs {
		if q.Root {
			roots++
		}
	}
	if roots != 3 {
		t.Errorf("root qdiscs = %d, want 3", roots)
	}
	if qs[1].Kind != "mq" || qs[1].Dev != "eth0" || !qs[1].Root {
		t.Errorf("eth0 root = %+v", qs[1])
	}
	if qs[2].Parent != ":1" || qs[2].Root {
		t.Errorf("child = %+v", qs[2])
	}
}
