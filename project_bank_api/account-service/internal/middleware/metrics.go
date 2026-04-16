package middleware

import (
	"account-service/pkg/telemetry"
	"net/http"
	"strconv"
	"time"
)

func PrometheusMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			h.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(sw, r)

		duration := time.Since(start).Seconds()
		statusStr := strconv.Itoa(sw.statusCode)

		telemetry.HttpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()
		telemetry.HttpRequestDuration.WithLabelValues(r.Method, r.URL.Path, statusStr).Observe(duration)
	}
}
