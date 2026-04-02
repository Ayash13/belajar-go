package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// ============================================================================
// DEFINISI METRICS — Mendaftarkan metrics yang akan di-expose ke Prometheus
// ============================================================================

var (
	// Counter: angka yang hanya naik (total request)
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_api_http_requests_total",
			Help: "Total jumlah HTTP request yang diterima",
		},
		[]string{"method", "path", "status"}, // label untuk filter
	)

	// Gauge: angka yang bisa naik-turun (active connections)
	activeConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "bank_api_active_connections",
			Help: "Jumlah koneksi HTTP yang sedang aktif",
		},
	)

	// Histogram: distribusi response time dalam bucket
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bank_api_http_duration_seconds",
			Help:    "Durasi HTTP request dalam detik",
			Buckets: prometheus.DefBuckets, // default: .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
		},
		[]string{"method", "path"},
	)

	// Counter: total transfer
	transferTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bank_api_transfers_total",
			Help: "Total jumlah transfer yang diproses",
		},
	)

	// Counter: total amount yang di-transfer (dalam rupiah)
	transferAmountTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "bank_api_transfer_amount_total",
			Help: "Total nominal transfer yang diproses (dalam Rupiah)",
		},
	)
)

func init() {
	// Mendaftarkan semua metrics ke Prometheus default registry
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(activeConnections)
	prometheus.MustRegister(httpDuration)
	prometheus.MustRegister(transferTotal)
	prometheus.MustRegister(transferAmountTotal)
}

// ============================================================================
// MIDDLEWARE — Otomatis mencatat metrics untuk setiap request
// ============================================================================

func metricsMiddleware(next http.HandlerFunc, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Naikkan gauge: ada koneksi baru yang aktif
		activeConnections.Inc()
		defer activeConnections.Dec() // Turunkan setelah selesai

		// Simulasi latency yang bervariasi
		delay := time.Duration(rand.Intn(200)) * time.Millisecond
		time.Sleep(delay)

		// Panggil handler yang sebenarnya
		next(w, r)

		// Catat metrics setelah request selesai
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(r.Method, path, "success").Inc()
		httpDuration.WithLabelValues(r.Method, path).Observe(duration)
	}
}

// ============================================================================
// HANDLERS — Endpoint untuk operasi banking
// ============================================================================

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accounts := []map[string]interface{}{
		{"id": "acc-001", "holder": "Ayash", "balance": 500000},
		{"id": "acc-002", "holder": "Budi", "balance": 300000},
		{"id": "acc-003", "holder": "Citra", "balance": 750000},
	}

	json.NewEncoder(w).Encode(response{Status: "success", Data: accounts})
}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	amount := 50000 + rand.Intn(450000) // random 50k - 500k

	// Catat metrics transfer
	transferTotal.Inc()
	transferAmountTotal.Add(float64(amount))

	json.NewEncoder(w).Encode(response{
		Status:  "success",
		Message: fmt.Sprintf("Transfer Rp %d berhasil diproses", amount),
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Status: "ok", Message: "Service is healthy"})
}

// ============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║   PROMETHEUS METRICS                             ║")
	fmt.Println("║   Instrumentasi HTTP Server dengan Prometheus    ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")

	mux := http.NewServeMux()

	// Endpoint bisnis — dibungkus metricsMiddleware
	mux.HandleFunc("GET /accounts", metricsMiddleware(handleGetAccounts, "/accounts"))
	mux.HandleFunc("POST /transfer", metricsMiddleware(handleTransfer, "/transfer"))
	mux.HandleFunc("GET /health", metricsMiddleware(handleHealth, "/health"))

	// Endpoint khusus Prometheus — di-scrape oleh Prometheus setiap 15 detik
	mux.Handle("GET /metrics", promhttp.Handler())

	fmt.Println("\n📊 Endpoints:")
	fmt.Println("  GET  http://localhost:8080/accounts  → data akun")
	fmt.Println("  POST http://localhost:8080/transfer   → simulasi transfer")
	fmt.Println("  GET  http://localhost:8080/health     → health check")
	fmt.Println("  GET  http://localhost:8080/metrics    → Prometheus metrics")
	fmt.Println("\n🚀 Server listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
