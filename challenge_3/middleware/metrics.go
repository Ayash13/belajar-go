package middleware

import (
	"belajar-go/challenge_3/telemetry"
	"net/http"
	"strconv"
	"time"
)

// PrometheusMiddleware merekam jumlah request masuk dan latensinya untuk disajikan ke sistem metrics prometheus
func PrometheusMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Jangan record metrik untuk endpoint /metrics itu sendiri agar tidak mengotori data
		if r.URL.Path == "/metrics" {
			h.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(sw, r)

		duration := time.Since(start).Seconds()

		// status dari integer (200), dikonversi jadi string ("200") untuk label metric
		statusStr := strconv.Itoa(sw.statusCode)

		// 1. Catat Total Request (Counter)
		telemetry.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()

		// 2. Catat Durasi Eksekusi (Histogram)
		telemetry.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, statusStr).Observe(duration)
	}
}
