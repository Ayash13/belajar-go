package middleware

import (
	"net"
	"net/http"
	"time"
	"transaction-service/config"
	"transaction-service/pkg/telemetry"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type statusWriter struct {
	http.ResponseWriter
	statusCode int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.statusCode = code
	sw.ResponseWriter.WriteHeader(code)
}

func RequestLogger(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/metrics" {
			h.ServeHTTP(w, r)
			return
		}
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(sw, r)
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		traceID := "none"
		if telemetry.Tracer != nil {
			traceID = trace.SpanContextFromContext(r.Context()).TraceID().String()
		}
		config.Log.Info("Request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.Int("status", sw.statusCode),
			zap.Duration("latency", time.Since(start)),
			zap.String("ip", ip),
			zap.String("trace_id", traceID),
		)
	}
}
