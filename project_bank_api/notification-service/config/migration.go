package config

import (
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

func RunMigrations(db *sqlx.DB) {
	migrationSQL := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS notifications (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		event_type VARCHAR(100) NOT NULL,
		reference_no VARCHAR(255) NOT NULL,
		account_no VARCHAR(50) NOT NULL,
		payload JSONB NOT NULL,
		callback_url TEXT NOT NULL,
		status VARCHAR(20) NOT NULL DEFAULT 'PENDING',
		retry_count INT NOT NULL DEFAULT 0,
		created_at TIMESTAMP NOT NULL DEFAULT NOW(),
		sent_at TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_notifications_reference ON notifications(reference_no);
	CREATE INDEX IF NOT EXISTS idx_notifications_account ON notifications(account_no);
	CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
	`

	if _, err := db.Exec(migrationSQL); err != nil {
		Log.Fatal("Failed to run migrations", zap.Error(err))
	}

	Log.Info("Database migration completed")
}
