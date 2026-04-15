package database

import (
	"belajar-go/challenge_3/logger"
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func ConnectDB() *sqlx.DB {
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn("Error loading .env file, falling back to environment variables")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || dbname == "" {
		logger.Log.Fatal("Database environment variables are missing")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}

	createTableQuery := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS accounts (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		account_holder VARCHAR(255) NOT NULL,
		balance NUMERIC(15,2) NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS transactions (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		partner_reference_no VARCHAR(255) UNIQUE NOT NULL,
		from_account_id UUID NOT NULL REFERENCES accounts(id),
		to_account_id UUID NOT NULL REFERENCES accounts(id),
		amount NUMERIC(15,2) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT NOW()
	);
	`
	if _, err = db.Exec(createTableQuery); err != nil {
		logger.Log.Fatal("Failed to create tables", zap.Error(err))
	}

	logger.Log.Info("Database connected successfully",
		zap.String("host", host),
		zap.String("dbname", dbname),
	)

	return db
}
