package middleware

import (
	"net/http"
	"strconv"
	"time"
	"transaction-service/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

func OpenTelemetryMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			h.ServeHTTP(w, r)
			return
		}
		if telemetry.Tracer == nil {
			h.ServeHTTP(w, r)
			return
		}
		ctx, span := telemetry.Tracer.Start(r.Context(), "HTTP "+r.Method+" "+r.URL.Path,
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.path", r.URL.Path),
				attribute.String("http.client_ip", r.RemoteAddr),
			),
		)
		defer span.End()
		r = r.WithContext(ctx)
		w.Header().Set("X-Trace-Id", span.SpanContext().TraceID().String())
		sw2 := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(sw2, r)
		span.SetAttributes(attribute.Int("http.status_code", sw2.statusCode))
	}
}
