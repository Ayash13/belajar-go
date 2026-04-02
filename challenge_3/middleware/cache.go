package middleware

import (
	"belajar-go/challenge_3/logger"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// responseRecorder untuk menangkap response body sebelum dikirim ke client agar bisa di-cache
type responseRecorder struct {
	http.ResponseWriter
	status int
	body   []byte
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body = append(r.body, b...)
	return r.ResponseWriter.Write(b)
}

func (r *responseRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

// Cache menyimpan response GET route ke Redis
// Naming Convention: bank_api:cache:url_path:{path}
func Cache(rdb *redis.Client, ttl time.Duration, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rdb == nil || r.Method != http.MethodGet {
			h.ServeHTTP(w, r)
			return
		}

		cacheKey := "bank_api:cache:url_path:" + r.URL.Path
		cachedData, err := rdb.Get(r.Context(), cacheKey).Result()
		if err == nil {
			logger.Log.Debug("Cache hit",
				zap.String("key", cacheKey),
			)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cachedData))
			return
		}

		rec := &responseRecorder{ResponseWriter: w, status: http.StatusOK}
		h.ServeHTTP(rec, r)

		if rec.status == http.StatusOK {
			if err := rdb.Set(r.Context(), cacheKey, rec.body, ttl).Err(); err != nil {
				logger.Log.Error("Cache set error",
					zap.String("key", cacheKey),
					zap.Error(err),
				)
			}
		}
	}
}
