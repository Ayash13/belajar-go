package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// JSONContentType menambahkan header Content-Type application/json ke setiap response
func JSONContentType(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	}
}

// HandleNotFound memberikan respons JSON khusus jika rute tidak terdaftar
func HandleNotFound(mux *http.ServeMux) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, pattern := mux.Handler(r)
		if pattern == "" {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"status":"error","message":"Route not found"}`))
			return
		}
		mux.ServeHTTP(w, r)
	}
}

// ApplyGlobalMiddlewares merangkum seluruh chain middleware di satu tempat
func ApplyGlobalMiddlewares(rdb *redis.Client, mux *http.ServeMux) http.Handler {
	var h http.Handler = HandleNotFound(mux)

	h = http.TimeoutHandler(h, 10*time.Second, `{"status":"error","message":"Server Timeout! Took too long to respond."}`)
	h = Idempotency(rdb, h)
	h = RateLimit(rdb, "global", 10, 5*time.Second, h)
	h = JSONContentType(h)
	h = RequestLogger(h)
	h = PrometheusMiddleware(h)
	h = OpenTelemetryMiddleware(h)

	return h
}
