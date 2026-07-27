package correlate

import (
	"testing"
	"time"

	"github.com/bgrewell/ifscope/internal/model"
)

func TestStatsRates(t *testing.T) {
	before := []model.InterfaceStats{{
		Name: "eth0", RxBytes: 1000, RxPackets: 10, RxErrors: 1, RxDropped: 2,
		TxBytes: 2000, TxPackets: 20, TxErrors: 2, TxDropped: 4,
	}}
	after := []model.InterfaceStats{{
		Name: "eth0", RxBytes: 3000, RxPackets: 30, RxErrors: 3, RxDropped: 6,
		TxBytes: 6000, TxPackets: 60, TxErrors: 6, TxDropped: 12,
		LinkSpeedBps: 10_000,
	}}

	got := StatsRates(before, after, 2*time.Second)
	rate := got[0].Rates
	if rate == nil {
		t.Fatal("rates are nil")
	}
	if rate.RxBitsPerSec != 8000 || rate.TxBitsPerSec != 16000 {
		t.Fatalf("bit rates = rx %v tx %v", rate.RxBitsPerSec, rate.TxBitsPerSec)
	}
	if rate.RxPacketsSec != 10 || rate.TxPacketsSec != 20 {
		t.Fatalf("packet rates = rx %v tx %v", rate.RxPacketsSec, rate.TxPacketsSec)
	}
	if rate.RxErrorsSec != 1 || rate.TxDroppedSec != 4 {
		t.Fatalf("error/drop rates = rx errors %v tx drops %v", rate.RxErrorsSec, rate.TxDroppedSec)
	}
	if rate.RxUtilization != 80 || rate.TxUtilization != 160 {
		t.Fatalf("utilization = rx %v tx %v", rate.RxUtilization, rate.TxUtilization)
	}
}

func TestStatsRatesMarksNewAndResetInterfaces(t *testing.T) {
	before := []model.InterfaceStats{{Name: "reset0", RxBytes: 100}}
	after := []model.InterfaceStats{{Name: "reset0", RxBytes: 10}, {Name: "new0", RxBytes: 20}}

	got := StatsRates(before, after, time.Second)
	if got[0].Rates != nil || got[0].RateStatus != RateStatusCounterReset {
		t.Fatalf("reset interface = %+v", got[0])
	}
	if got[1].Rates != nil || got[1].RateStatus != RateStatusNew {
		t.Fatalf("new interface = %+v", got[1])
	}
}

func TestStatsRatesMarksRecreatedInterface(t *testing.T) {
	got := StatsRates(
		[]model.InterfaceStats{{ID: 2, Name: "eth0", RxBytes: 10}},
		[]model.InterfaceStats{{ID: 9, Name: "eth0", RxBytes: 20}},
		time.Second,
	)
	if got[0].Rates != nil || got[0].RateStatus != RateStatusRecreated {
		t.Fatalf("recreated interface = %+v", got[0])
	}
}

func TestStatsRatesHandlesCounterWrap(t *testing.T) {
	got := StatsRates(
		[]model.InterfaceStats{{Name: "eth0", RxBytes: 0xfffffff0}},
		[]model.InterfaceStats{{Name: "eth0", RxBytes: 0x10}},
		time.Second,
	)
	if got[0].Rates == nil || got[0].Rates.RxBitsPerSec != 32*8 {
		t.Fatalf("wrapped rates = %+v", got[0])
	}
}

func TestStatsRatesRejectsNonPositiveElapsedTime(t *testing.T) {
	got := StatsRates(
		[]model.InterfaceStats{{Name: "eth0", RxBytes: 1}},
		[]model.InterfaceStats{{Name: "eth0", RxBytes: 2}},
		0,
	)
	if got[0].Rates != nil {
		t.Fatalf("rates = %+v, want nil", got[0].Rates)
	}
}

func TestSortStats(t *testing.T) {
	stats := []model.InterfaceStats{
		{Name: "slow", Rates: &model.InterfaceStatsRate{RxBitsPerSec: 1}},
		{Name: "fast", Rates: &model.InterfaceStatsRate{RxBitsPerSec: 10}},
	}
	SortStats(stats, "rx")
	if stats[0].Name != "fast" {
		t.Fatalf("first = %q, want fast", stats[0].Name)
	}
}
