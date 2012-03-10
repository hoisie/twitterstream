package httpstream

import (
	"sync"
	"time"
)

var counterMu sync.Mutex
var streamCounters = make(map[string]uint64)
var streamGauges = make(map[string]float64)
var counterHistory = make([]map[string]uint64, 0)

func init() {
	RunMetricsHeartbeat()
}

// get this metrics counters for the last 10 tests
func MetricsHistory() []map[string]uint64 {
	// hopefully this returns a copy?
	return counterHistory
}


func incrCounter(name string) {
	updateCounter(name, 1)
}

func updateGauge(name string, value float64) {
	counterMu.Lock()
	streamGauges[name] += value
	counterMu.Unlock()
}

func updateCounter(name string, value uint64) {
	counterMu.Lock()
	streamCounters[name] += value
	counterMu.Unlock()
}

// gets metrics, and replaces with new zeros, keep a history of last 10 samples
// so we can calculate throughput rolling avgs
func MetricsSnapshot() (counters map[string]uint64, gauges map[string]float64) {
	counterMu.Lock()
	counters = streamCounters
	gauges = streamGauges
	streamCounters = make(map[string]uint64)
	streamGauges = make(map[string]float64)
	if len(counterHistory) > 10 {
		counterHistory = counterHistory[1:]
	}
	counterHistory = append(counterHistory, counters)
	counterMu.Unlock()
	return
}

func RunMetricsHeartbeat() {

	// lets poll back to take metrics snapshots
	timer := time.NewTicker(time.Second * 60)

	go func() {
		for _ = range timer.C {
			MetricsSnapshot()
		}
	}()
}
