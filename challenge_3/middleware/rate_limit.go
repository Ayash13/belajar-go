package middleware

import (
	"belajar-go/challenge_3/logger"
	"net"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RateLimit menggunakan Redis untuk membatasi jumlah request per IP dalam satu jendela waktu (window)
// Naming Convention: bank_api:rate_limit:{route_prefix}:ip:{ip}
func RateLimit(rdb *redis.Client, routePrefix string, maxReq int, window time.Duration, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rdb == nil {
			h.ServeHTTP(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		key := "bank_api:rate_limit:" + routePrefix + ":ip:" + ip

		count, err := rdb.Incr(r.Context(), key).Result()
		if err != nil {
			logger.Log.Error("Rate limit Redis error",
				zap.String("key", key),
				zap.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"status":"error","message":"Redis server error handling rate limit"}`))
			return
		}

		if count == 1 {
			rdb.Expire(r.Context(), key, window)
		}

		if count > int64(maxReq) {
			logger.Log.Warn("Rate limit exceeded",
				zap.String("ip", ip),
				zap.String("route_prefix", routePrefix),
				zap.Int64("count", count),
				zap.Int("max", maxReq),
			)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"status":"error","message":"Too Many Requests! Rate Limit Exceeded for this IP"}`))
			return
		}

		h.ServeHTTP(w, r)
	}
}
