package correlate

import (
	"math"
	"sort"
	"time"

	"github.com/bgrewell/ifscope/internal/model"
)

const (
	RateStatusNew          = "new_interface"
	RateStatusRecreated    = "recreated_interface"
	RateStatusCounterReset = "counter_reset"
)

// StatsRates derives per-second rates by matching current interface counters
// to a previous sample. New interfaces and reset/decreased counters are marked
// explicitly rather than producing misleading underflow rates.
func StatsRates(previous, current []model.InterfaceStats, elapsed time.Duration) []model.InterfaceStats {
	out := append([]model.InterfaceStats(nil), current...)
	if elapsed <= 0 {
		return out
	}

	byName := make(map[string]model.InterfaceStats, len(previous))
	for _, stat := range previous {
		byName[stat.Name] = stat
	}
	seconds := elapsed.Seconds()
	for i := range out {
		before, ok := byName[out[i].Name]
		if !ok {
			out[i].RateStatus = RateStatusNew
			continue
		}
		if before.ID != 0 && out[i].ID != 0 && before.ID != out[i].ID {
			out[i].RateStatus = RateStatusRecreated
			continue
		}
		delta, ok := counterDeltas(before, out[i])
		if !ok {
			out[i].RateStatus = RateStatusCounterReset
			continue
		}

		rate := &model.InterfaceStatsRate{
			SampleSeconds: seconds,
			RxBitsPerSec:  float64(delta.RxBytes) * 8 / seconds,
			RxPacketsSec:  float64(delta.RxPackets) / seconds,
			RxErrorsSec:   float64(delta.RxErrors) / seconds,
			RxDroppedSec:  float64(delta.RxDropped) / seconds,
			TxBitsPerSec:  float64(delta.TxBytes) * 8 / seconds,
			TxPacketsSec:  float64(delta.TxPackets) / seconds,
			TxErrorsSec:   float64(delta.TxErrors) / seconds,
			TxDroppedSec:  float64(delta.TxDropped) / seconds,
		}
		if out[i].LinkSpeedBps > 0 {
			rate.RxUtilization = rate.RxBitsPerSec / float64(out[i].LinkSpeedBps) * 100
			rate.TxUtilization = rate.TxBitsPerSec / float64(out[i].LinkSpeedBps) * 100
		}
		out[i].Rates = rate
	}
	return out
}

func counterDeltas(a, b model.InterfaceStats) (model.InterfaceStats, bool) {
	var delta model.InterfaceStats
	pairs := []struct {
		before, after uint64
		target        *uint64
	}{
		{a.RxBytes, b.RxBytes, &delta.RxBytes},
		{a.RxPackets, b.RxPackets, &delta.RxPackets},
		{a.RxErrors, b.RxErrors, &delta.RxErrors},
		{a.RxDropped, b.RxDropped, &delta.RxDropped},
		{a.TxBytes, b.TxBytes, &delta.TxBytes},
		{a.TxPackets, b.TxPackets, &delta.TxPackets},
		{a.TxErrors, b.TxErrors, &delta.TxErrors},
		{a.TxDropped, b.TxDropped, &delta.TxDropped},
		{a.Collisions, b.Collisions, &delta.Collisions},
	}
	for _, pair := range pairs {
		value, ok := counterDelta(pair.before, pair.after)
		if !ok {
			return model.InterfaceStats{}, false
		}
		*pair.target = value
	}
	return delta, true
}

// counterDelta accepts monotonic growth and wraps close to the 32- or 64-bit
// boundary. Other decreases are treated as a counter reset.
func counterDelta(before, after uint64) (uint64, bool) {
	if after >= before {
		return after - before, true
	}
	if before >= 0xf0000000 && before <= math.MaxUint32 {
		return math.MaxUint32 - before + 1 + after, true
	}
	if before >= math.MaxUint64-0xffff {
		return math.MaxUint64 - before + 1 + after, true
	}
	return 0, false
}

// SortStats orders stats by a rate metric descending, or by name ascending.
func SortStats(stats []model.InterfaceStats, by string) {
	sort.SliceStable(stats, func(i, j int) bool {
		if by == "name" {
			return stats[i].Name < stats[j].Name
		}
		iv, jv := statsMetric(stats[i], by), statsMetric(stats[j], by)
		if iv == jv {
			return stats[i].Name < stats[j].Name
		}
		return iv > jv
	})
}

func statsMetric(stat model.InterfaceStats, by string) float64 {
	if stat.Rates == nil {
		return 0
	}
	switch by {
	case "rx":
		return stat.Rates.RxBitsPerSec
	case "tx":
		return stat.Rates.TxBitsPerSec
	case "errors":
		return stat.Rates.RxErrorsSec + stat.Rates.TxErrorsSec
	case "drops":
		return stat.Rates.RxDroppedSec + stat.Rates.TxDroppedSec
	case "utilization":
		return max(stat.Rates.RxUtilization, stat.Rates.TxUtilization)
	default:
		return 0
	}
}
