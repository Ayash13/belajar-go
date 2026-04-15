package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load("../.env")
	}
}

func ConnectDB() *sqlx.DB {
	host := os.Getenv("SNAP_DB_HOST")
	port := os.Getenv("SNAP_DB_PORT")
	user := os.Getenv("SNAP_DB_USER")
	password := os.Getenv("SNAP_DB_PASSWORD")
	dbname := os.Getenv("SNAP_DB_NAME")
	sslmode := os.Getenv("SNAP_DB_SSLMODE")

	if host == "" || user == "" || dbname == "" {
		Log.Fatal("Database environment variables are missing (SNAP_DB_*)")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	Log.Info("Database connected successfully",
		zap.String("host", host),
		zap.String("dbname", dbname),
	)

	return db
}

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
		Log.Fatal("Failed to connect to Redis", zap.Error(err))
	}

	Log.Info("Redis connected successfully",
		zap.String("addr", redisHost+":"+redisPort),
	)
	return rdb
}
