package parse

import "testing"

func TestIPRules(t *testing.T) {
	// A typical source-based routing setup: from a subnet, look up a custom
	// table; plus the three default rules.
	data := []byte(`[
	  {"priority":0,"src":"all","table":"local"},
	  {"priority":100,"src":"192.168.9.0","srclen":24,"table":"50"},
	  {"priority":101,"src":"all","fwmark":"0x1","table":"51"},
	  {"priority":32766,"src":"all","table":"main"}
	]`)
	rules, err := IPRules(data)
	if err != nil {
		t.Fatal(err)
	}
	if len(rules) != 4 {
		t.Fatalf("rules = %d, want 4", len(rules))
	}

	src := rules[1]
	if src.Priority != 100 || src.From != "192.168.9.0/24" || src.Table != "50" {
		t.Errorf("source-based rule = %+v", src)
	}
	if src.Family != "inet" {
		t.Errorf("family = %q", src.Family)
	}

	fw := rules[2]
	if fw.FWMark != "0x1" || fw.Table != "51" {
		t.Errorf("fwmark rule = %+v", fw)
	}

	if rules[0].From != "all" {
		t.Errorf("from = %q, want all", rules[0].From)
	}
}

func TestIPRulesIPv6Family(t *testing.T) {
	data := []byte(`[{"priority":100,"src":"2001:db8::","srclen":32,"table":"80"}]`)
	rules, _ := IPRules(data)
	if rules[0].Family != "inet6" {
		t.Errorf("family = %q, want inet6", rules[0].Family)
	}
	if rules[0].From != "2001:db8::/32" {
		t.Errorf("from = %q", rules[0].From)
	}
}
