package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type Metric struct {
	Timestamp   time.Time
	LatencyMs   float64
	MessagesSec int64
}

var (
	mu           sync.Mutex
	latencySum   int64
	latencyCount int64
	sentCount    int64

	samples    []Metric
	maxSamples = 300
)

func recordLatency(d time.Duration) {
	atomic.AddInt64(&latencySum, d.Microseconds())
	atomic.AddInt64(&latencyCount, 1)
}

func recordSend() {
	atomic.AddInt64(&sentCount, 1)
}

func startMetricsServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "metrics.html")
	})
	http.HandleFunc("/stats", serveStats)

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			sum := atomic.SwapInt64(&latencySum, 0)
			cnt := atomic.SwapInt64(&latencyCount, 0)
			sent := atomic.SwapInt64(&sentCount, 0)

			var avg float64
			if cnt > 0 {
				avg = float64(sum) / float64(cnt) / 1000.0
			}

			mu.Lock()
			samples = append(samples, Metric{
				Timestamp:   time.Now(),
				LatencyMs:   avg,
				MessagesSec: sent,
			})
			if len(samples) > maxSamples {
				samples = samples[len(samples)-maxSamples:]
			}
			mu.Unlock()
		}
	}()

	fmt.Println("Metrics dashboard at: http://localhost:9090")
	go http.ListenAndServe(":9090", nil)
}

func serveStats(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(samples)
}