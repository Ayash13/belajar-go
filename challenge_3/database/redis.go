package database

import (
	"belajar-go/challenge_3/logger"
	"context"
	"os"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisHost + ":" + redisPort,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		logger.Log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	logger.Log.Info("Redis connected successfully",
		zap.String("addr", redisHost+":"+redisPort),
	)
	return rdb
}
