package metrics

import (
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"runtime"
	"time"
)

type MetricsMap map[string]RouteMetrics

type RouteMetrics struct {
	Route          string `json:"route"`
	TotalRequests  int64  `json:"totalRequests"`
	DurationTotal  int64  `json:"durationTotal"`
	DurationAvg    int64  `json:"durationAvg"`
	DurationRecent int64  `json:"durationRecent"`
}

var metricsMap = MetricsMap{}

func PerformanceLoggingMiddleware(next http.Handler, logger zerolog.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		if _, ok := metricsMap[r.URL.String()]; !ok {
			newMetrics := RouteMetrics{
				Route:          r.URL.String(),
				TotalRequests:  1,
				DurationTotal:  int64(duration),
				DurationAvg:    int64(duration),
				DurationRecent: int64(duration),
			}
			logger.Info().Interface("metrics", newMetrics).Msg("Metrics recorded")
			metricsMap[r.URL.String()] = newMetrics
		} else {
			oldMetrics := metricsMap[r.URL.String()]
			newMetrics := RouteMetrics{
				Route:          oldMetrics.Route,
				TotalRequests:  oldMetrics.TotalRequests + 1,
				DurationTotal:  int64(duration) + oldMetrics.DurationTotal,
				DurationAvg:    int64(duration) + oldMetrics.DurationTotal/oldMetrics.TotalRequests + 1,
				DurationRecent: int64(duration),
			}
			logger.Info().Interface("new metrics", newMetrics).Msg("route exists " + r.URL.String())
			metricsMap[r.URL.String()] = newMetrics
		}
		logger.Info().Str("method", r.Method).Str("url", r.URL.String()).Str("duration", duration.String()).Msg("metrics")
	})
}

func HandleSendingMetrics(w http.ResponseWriter, r *http.Request, logger zerolog.Logger) {
	_, err := w.Write([]byte("<h1>Metrics for each route. Times are in nanoseconds.</h1>"))
	if err != nil {
		logger.Error().Err(err).Msg("Unable to write response")
	}
	for _, metrics := range metricsMap {
		_, err = w.Write([]byte(fmt.Sprintf("<h2>Route: %s</h2>", metrics.Route)))
		_, err = w.Write([]byte(fmt.Sprintf("<h3>Total Requests: %d</h3>", metrics.TotalRequests)))
		_, err = w.Write([]byte(fmt.Sprintf("<h3>Most Recent Request Duration: %d</h3>", metrics.DurationRecent)))
		_, err = w.Write([]byte(fmt.Sprintf("<h3>Avg Request Duration: %d</h3>", metrics.DurationAvg)))
		if err != nil {
			logger.Error().Err(err).Msg("Unable to encode metrics")
		}
		_, err = w.Write([]byte("\n"))
	}
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	_, err = w.Write([]byte("<h2>Memory Stats:</h2>"))
	_, err = w.Write([]byte(fmt.Sprintf("<h3>Allocated memory: %v</h3>", memStats.Alloc)))
	_, err = w.Write([]byte("<br>"))
	if err != nil {
		logger.Error().Err(err).Msg("Unable to write response")
	}
	_, err = w.Write([]byte(fmt.Sprintf("<h3>Heap allocated memory: %v</h3>", memStats.HeapAlloc)))
	_, err = w.Write([]byte("<br>"))
	if err != nil {
		logger.Error().Err(err).Msg("Unable to write response")
	}
	_, err = w.Write([]byte(fmt.Sprintf("<h3>System memory obtained from OS: %v</h3>", memStats.Sys)))
	_, err = w.Write([]byte("<br>"))
	if err != nil {
		logger.Error().Err(err).Msg("Unable to write response")
	}
	_, err = w.Write([]byte(fmt.Sprintf("<h3>Number of garbage collection cycles: %v</h3>", memStats.NumGC)))
	_, err = w.Write([]byte("<br>"))
	if err != nil {
		logger.Error().Err(err).Msg("Unable to write response")
	}
}
