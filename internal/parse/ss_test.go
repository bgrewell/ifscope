package parse

import "testing"

func TestSS(t *testing.T) {
	// Headerless `ss -tulpnH` lines, incl. IPv6, a wildcard, and a process col.
	data := []byte(`udp UNCONN 0 0 127.0.0.53%lo:53 0.0.0.0:* users:(("systemd-resolve",pid=1971729,fd=17))
tcp LISTEN 0 4096 *:22 *:*
tcp LISTEN 0 128 [::]:443 [::]:* users:(("nginx",pid=42,fd=6),("nginx",pid=43,fd=6))`)

	socks := SS(data)
	if len(socks) != 3 {
		t.Fatalf("sockets = %d, want 3", len(socks))
	}

	dns := socks[0]
	if dns.Proto != "udp" || dns.LocalAddr != "127.0.0.53%lo" || dns.LocalPort != "53" {
		t.Errorf("dns = %+v", dns)
	}
	if dns.Process != "systemd-resolve(1971729)" {
		t.Errorf("dns process = %q", dns.Process)
	}

	ssh := socks[1]
	if ssh.LocalAddr != "*" || ssh.LocalPort != "22" || ssh.Process != "" {
		t.Errorf("ssh = %+v", ssh)
	}

	https := socks[2]
	if https.LocalAddr != "[::]" || https.LocalPort != "443" {
		t.Errorf("https local = %+v", https)
	}
	if https.Process != "nginx(42),nginx(43)" {
		t.Errorf("https process = %q", https.Process)
	}
}

func TestSSEmpty(t *testing.T) {
	if got := SS(nil); len(got) != 0 {
		t.Errorf("empty = %v", got)
	}
}
