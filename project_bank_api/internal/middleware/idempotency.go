package middleware

import (
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func Idempotency(rdb *redis.Client, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			h.ServeHTTP(w, r)
			return
		}

		key := r.Header.Get("Idempotency-Key")
		if key != "" && rdb != nil {
			redisKey := "snap_api:idempotency:key:" + key
			isFirstTime, err := rdb.SetNX(r.Context(), redisKey, "started", 500*time.Millisecond).Result()
			if err != nil || !isFirstTime {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(`{"responseCode":"4090000","responseMessage":"Duplicate transaction rejected"}`))
				return
			}
		}
		h.ServeHTTP(w, r)
	}
}
