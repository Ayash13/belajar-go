package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

// Idempotency API untuk mutasi mencegah duplikasi
func Idempotency(rdb *redis.Client, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Abaikan GET requests karena sudah idempotent secara definisi
		if r.Method == http.MethodGet {
			h.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("Idempotency-Key")
		if key != "" && rdb != nil {
			// bank_api:idempotency:key:uuid-123-abc
			redisKey := "bank_api:idempotency:key:" + key
			isFirstTime, err := rdb.SetNX(r.Context(), redisKey, "started", 24*time.Hour).Result()
			if err != nil || !isFirstTime {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(`{"status":"error","message":"duplicate transaction rejected"}`))
				return
			}
		}
		h.ServeHTTP(w, r)
	}
}
