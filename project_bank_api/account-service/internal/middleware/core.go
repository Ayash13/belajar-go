package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func JSONContentType(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

func HandleNotFound(mux *http.ServeMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"responseCode":"4040000","responseMessage":"Service Not Found"}`))
			return
		}
		mux.ServeHTTP(w, r)
	}
}

func CorsMiddleware(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Authorization-Customer, X-TIMESTAMP, X-SIGNATURE, ORIGIN, X-PARTNER-ID, X-EXTERNAL-ID, X-IP-ADDRESS, X-DEVICE-ID, X-LATITUDE, X-LONGITUDE, CHANNEL-ID, Idempotency-Key")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	}
}

func ApplyGlobalMiddlewares(rdb *redis.Client, inner http.Handler) http.Handler {
	var h http.Handler = inner

	h = http.TimeoutHandler(h, 10*time.Second, `{"responseCode":"5040000","responseMessage":"Server Timeout"}`)
	h = Idempotency(rdb, h)
	h = RateLimit(rdb, "global", 10, 5*time.Second, h)
	h = JSONContentType(h)
	h = CorsMiddleware(h)
	h = RequestLogger(h)
	h = PrometheusMiddleware(h)
	h = OpenTelemetryMiddleware(h)

	return h
}
