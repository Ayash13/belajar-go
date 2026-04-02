package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync/atomic"
	"time"
)

// ============================================================================
// PILAR 1: LOGS — Catatan detail dari setiap kejadian dalam aplikasi
// ============================================================================

// simulateUnstructuredLog menunjukkan cara logging tradisional yang sulit di-parse mesin
func simulateUnstructuredLog() {
	fmt.Println("\n📝 ===== PILAR 1: LOGS =====")
	fmt.Println("\n--- Unstructured Log (Tradisional) ---")

	// Log biasa — mudah dibaca manusia, tapi sulit di-query oleh tools seperti Loki
	log.Println("Server started on port 8080")
	log.Println("User ayash logged in")
	log.Printf("Transfer Rp 500000 dari akun A ke akun B berhasil")
	log.Println("ERROR: Database connection timeout after 5s")
}

// simulateStructuredLog menunjukkan structured logging yang mudah di-parse mesin
func simulateStructuredLog() {
	fmt.Println("\n--- Structured Log (slog / Zap) ---")

	// Structured log — output dalam format JSON, mudah di-query oleh Loki/Elasticsearch
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("Server started",
		slog.String("port", "8080"),
		slog.String("env", "production"),
	)

	logger.Info("User login berhasil",
		slog.String("user_id", "usr-123"),
		slog.String("username", "ayash"),
		slog.String("ip", "192.168.1.10"),
	)

	logger.Info("Transfer berhasil",
		slog.String("from_account", "acc-001"),
		slog.String("to_account", "acc-002"),
		slog.Int("amount", 500000),
		slog.String("currency", "IDR"),
	)

	logger.Error("Database connection timeout",
		slog.String("host", "db.internal:5432"),
		slog.Duration("timeout", 5*time.Second),
		slog.Int("retry_count", 3),
	)
}

// ============================================================================
// PILAR 2: METRICS — Data numerik yang diukur dari waktu ke waktu
// ============================================================================

// SimpleMetrics menyimpan metrics sederhana (di production pakai Prometheus)
type SimpleMetrics struct {
	totalRequests  atomic.Int64
	totalErrors    atomic.Int64
	totalLatencyMs atomic.Int64
}

func simulateMetrics() {
	fmt.Println("\n📊 ===== PILAR 2: METRICS =====")

	metrics := &SimpleMetrics{}

	// Simulasi beberapa request masuk
	requests := []struct {
		path      string
		latencyMs int64
		isError   bool
	}{
		{"/accounts", 45, false},
		{"/accounts/1", 32, false},
		{"/transfer", 180, false},
		{"/accounts", 12, false}, // dari cache, lebih cepat
		{"/transfer", 0, true},   // error
		{"/accounts", 55, false},
	}

	for _, req := range requests {
		metrics.totalRequests.Add(1)
		if req.isError {
			metrics.totalErrors.Add(1)
		}
		metrics.totalLatencyMs.Add(req.latencyMs)
	}

	totalReq := metrics.totalRequests.Load()
	totalErr := metrics.totalErrors.Load()
	avgLatency := metrics.totalLatencyMs.Load() / totalReq

	// Di production, angka-angka ini di-expose ke Prometheus via /metrics endpoint
	fmt.Printf("\n--- Metrics Summary ---\n")
	fmt.Printf("  http_requests_total     = %d (counter)\n", totalReq)
	fmt.Printf("  http_errors_total       = %d (counter)\n", totalErr)
	fmt.Printf("  http_latency_avg_ms     = %d ms (gauge)\n", avgLatency)
	fmt.Printf("  error_rate              = %.1f%% (gauge)\n", float64(totalErr)/float64(totalReq)*100)
	fmt.Println()
	fmt.Println("  Di production, data ini di-scrape oleh Prometheus")
	fmt.Println("  dan divisualisasikan di Grafana dashboard.")
}

// ============================================================================
// PILAR 3: TRACES — Melacak perjalanan request dari awal hingga akhir
// ============================================================================

// Span merepresentasikan satu langkah operasi dalam sebuah trace
type Span struct {
	Name     string
	Duration time.Duration
}

// simulateTrace mendemonstrasikan bagaimana satu request di-trace melewati berbagai layer
func simulateTrace() {
	fmt.Println("\n🔍 ===== PILAR 3: TRACES (Distributed Tracing) =====")

	traceID := "trace-abc-123-xyz"
	fmt.Printf("\n--- Trace ID: %s ---\n", traceID)
	fmt.Println("  Melacak request: GET /accounts/42")
	fmt.Println()

	// Simulasi span-span dalam satu trace
	spans := []Span{
		{"middleware.RateLimit", 2 * time.Millisecond},
		{"middleware.Cache (MISS)", 1 * time.Millisecond},
		{"handler.GetAccountByID", 3 * time.Millisecond},
		{"service.GetAccountByID", 5 * time.Millisecond},
		{"repository.GetByID (SQL Query)", 180 * time.Millisecond},
		{"middleware.Cache (SET)", 2 * time.Millisecond},
	}

	totalDuration := time.Duration(0)
	for _, span := range spans {
		totalDuration += span.Duration
	}

	// Visualisasi waterfall trace
	fmt.Println("  ┌─ Waterfall View ─────────────────────────────────────────┐")
	for _, span := range spans {
		barLen := int(float64(span.Duration) / float64(totalDuration) * 40)
		if barLen < 1 {
			barLen = 1
		}
		bar := ""
		for j := 0; j < barLen; j++ {
			bar += "█"
		}
		fmt.Printf("  │ %-38s %s %v\n", span.Name, bar, span.Duration)
	}
	fmt.Println("  └──────────────────────────────────────────────────────────┘")
	fmt.Printf("  Total Duration: %v\n", totalDuration)
	fmt.Println()
	fmt.Println("  Dari trace ini, kita bisa lihat bahwa bottleneck ada di")
	fmt.Println("  repository.GetByID (SQL Query) = 180ms dari total 193ms.")
	fmt.Println("  Di production, data ini dikirim ke Tempo/Jaeger via OpenTelemetry.")
}

// ============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║   OBSERVABILITY FUNDAMENTALS                    ║")
	fmt.Println("║   3 Pilar: Logs, Metrics, Traces                ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")

	// Pilar 1: Logs
	simulateUnstructuredLog()
	simulateStructuredLog()

	// Pilar 2: Metrics
	simulateMetrics()

	// Pilar 3: Traces
	simulateTrace()

	fmt.Println("\n✅ Ketiga pilar ini bekerja sama untuk memberikan")
	fmt.Println("   full visibility ke dalam sistem kita.")
}
