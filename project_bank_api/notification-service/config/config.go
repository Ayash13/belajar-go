package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		_ = godotenv.Load("../.env")
	}
}

func ConnectDB() *sqlx.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || dbname == "" {
		Log.Fatal("Database environment variables are missing (DB_*)")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	Log.Info("Database connected",
		zap.String("host", host),
		zap.String("dbname", dbname),
	)

	return db
}
