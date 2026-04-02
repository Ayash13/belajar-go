package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ============================================================================
// DEVELOPMENT LOGGER — Human-readable, berwarna, cocok saat development
// ============================================================================

func demoDevelopmentLogger() {
	fmt.Println("\n🛠️  ===== DEVELOPMENT LOGGER =====")

	// NewDevelopment() membuat logger dengan format human-readable
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flush buffer sebelum program selesai

	logger.Debug("Debug message — hanya tampil di development",
		zap.String("module", "auth"),
	)

	logger.Info("User berhasil login",
		zap.String("user_id", "usr-123"),
		zap.String("username", "ayash"),
		zap.String("ip", "192.168.1.10"),
	)

	logger.Warn("Response time lambat",
		zap.String("endpoint", "/accounts"),
		zap.Duration("latency", 2500*time.Millisecond),
		zap.Int("threshold_ms", 1000),
	)

	logger.Error("Gagal koneksi ke database",
		zap.String("host", "db.internal:5432"),
		zap.Int("retry_count", 3),
		zap.Duration("timeout", 5*time.Second),
	)
}

// ============================================================================
// PRODUCTION LOGGER — Output JSON, cocok di-ingest oleh Loki/Elasticsearch
// ============================================================================

func demoProductionLogger() {
	fmt.Println("\n🏭 ===== PRODUCTION LOGGER (JSON) =====")

	// NewProduction() membuat logger dengan output JSON murni
	// Level minimum default: Info (Debug tidak tampil)
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Debug("Ini TIDAK akan tampil di production") // di-skip karena level minimum = Info

	logger.Info("Transfer berhasil diproses",
		zap.String("trace_id", "trace-abc-123"),
		zap.String("from_account", "acc-001"),
		zap.String("to_account", "acc-002"),
		zap.Int("amount", 500000),
		zap.String("currency", "IDR"),
		zap.Duration("processing_time", 180*time.Millisecond),
	)

	logger.Warn("Rate limit hampir tercapai",
		zap.String("ip", "192.168.1.50"),
		zap.Int("current_count", 8),
		zap.Int("max_allowed", 10),
		zap.Duration("window", 5*time.Second),
	)

	logger.Error("Transfer gagal",
		zap.String("trace_id", "trace-def-456"),
		zap.String("from_account", "acc-003"),
		zap.String("error", "insufficient balance"),
		zap.Int("requested_amount", 1000000),
		zap.Int("available_balance", 500000),
	)
}

// ============================================================================
// SUGARED LOGGER — Lebih fleksibel, syntax mirip fmt.Printf
// ============================================================================

func demoSugaredLogger() {
	fmt.Println("\n🍬 ===== SUGARED LOGGER =====")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	// Sugar() mengonversi Logger ke SugaredLogger
	sugar := logger.Sugar()

	// Infow = Info with key-value pairs (lebih fleksibel dari zap.String, zap.Int)
	sugar.Infow("Cache hit",
		"key", "bank_api:cache:url_path:/accounts",
		"ttl_remaining_sec", 245,
		"source", "redis",
	)

	// Infof = Info with format string (mirip fmt.Sprintf)
	sugar.Infof("Akun %s berhasil dibuat dengan saldo awal Rp %d", "acc-099", 100000)

	// Errorw = Error with key-value pairs
	sugar.Errorw("Idempotency key sudah digunakan",
		"key", "pay-abc-123",
		"status", "rejected",
		"original_time", "2026-04-01T09:00:00Z",
	)
}

// ============================================================================
// CUSTOM LOGGER — Konfigurasi manual untuk kebutuhan khusus
// ============================================================================

func demoCustomLogger() {
	fmt.Println("\n⚙️  ===== CUSTOM LOGGER =====")

	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development: false,
		Encoding:    "json", // "json" atau "console"
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	defer logger.Sync()

	logger.Info("Custom logger dengan format timestamp ISO8601",
		zap.String("service", "bank_api"),
		zap.String("env", "production"),
		zap.String("version", "1.0.0"),
	)

	logger.Info("Request diterima",
		zap.String("method", "POST"),
		zap.String("path", "/transfer"),
		zap.String("trace_id", "trace-xyz-789"),
		zap.String("client_ip", "10.0.0.5"),
	)
}

// ============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║   STRUCTURED LOGGING WITH ZAP                   ║")
	fmt.Println("║   High-performance logging dari Uber             ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")

	demoDevelopmentLogger()
	demoProductionLogger()
	demoSugaredLogger()
	demoCustomLogger()

	fmt.Println("\n✅ Zap menghasilkan structured log (JSON) yang siap")
	fmt.Println("   di-ingest oleh Loki, Elasticsearch, atau CloudWatch.")
}
