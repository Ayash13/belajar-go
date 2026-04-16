package config

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RunMigrations(db *sqlx.DB) {
	migrationSQL := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS accounts (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		account_no VARCHAR(50) UNIQUE NOT NULL,
		customer_id VARCHAR(100) NOT NULL,
		name VARCHAR(255) NOT NULL,
		email VARCHAR(255) NOT NULL,
		phone_no VARCHAR(50),
		balance NUMERIC(15,2) NOT NULL DEFAULT 0,
		currency VARCHAR(10) NOT NULL DEFAULT 'IDR',
		partner_reference_no VARCHAR(255) UNIQUE NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMP NOT NULL DEFAULT NOW()
	);

	CREATE INDEX IF NOT EXISTS idx_accounts_account_no ON accounts(account_no);
	CREATE INDEX IF NOT EXISTS idx_accounts_partner_ref ON accounts(partner_reference_no);
	`

	if _, err := db.Exec(migrationSQL); err != nil {
		Log.Fatal("Failed to run migrations", zap.Error(err))
	}

	Log.Info("Database migration completed")
}
