package middleware

import (
	"account-service/pkg/telemetry"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

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
				attribute.String("http.user_agent", r.UserAgent()),
			),
		)
		defer span.End()

		r = r.WithContext(ctx)

		traceID := span.SpanContext().TraceID().String()
		w.Header().Set("X-Trace-Id", traceID)

		sw := &statusWriter{ResponseWriter: w, statusCode: http.StatusOK}
		h.ServeHTTP(sw, r)

		span.SetAttributes(attribute.Int("http.status_code", sw.statusCode))
		if sw.statusCode >= 500 {
			span.SetAttributes(attribute.Bool("error", true))
		}
	}
}
