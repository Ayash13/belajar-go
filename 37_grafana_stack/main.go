package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ============================================================================
// PROMETHEUS METRICS
// ============================================================================

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "bank_api_http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"method", "path", "status_code"},
	)

	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "bank_api_http_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpDuration)
}

// ============================================================================
// OPENTELEMETRY TRACER
// ============================================================================

var tracer = otel.Tracer("bank-api")

func initTracer() (*sdktrace.TracerProvider, error) {
	tempoHost := os.Getenv("TEMPO_HOST")
	if tempoHost == "" {
		tempoHost = "localhost"
	}

	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint(tempoHost+":4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, _ := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("bank-api"),
			semconv.ServiceVersionKey.String("1.0.0"),
		),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

// ============================================================================
// ZAP LOGGER — Output JSON ke stdout, dipick up Loki via Promtail/Alloy
// ============================================================================

func initLogger() *zap.Logger {
	config := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.InfoLevel),
		Development: false,
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "timestamp",
			LevelKey:       "level",
			MessageKey:     "message",
			CallerKey:      "caller",
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeDuration: zapcore.MillisDurationEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, _ := config.Build()
	return logger
}

// ============================================================================
// STATUS RECORDER — Menangkap status code untuk metrics
// ============================================================================

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// ============================================================================
// MIDDLEWARE — Menggabungkan Logging, Tracing, dan Metrics
// ============================================================================

func observabilityMiddleware(logger *zap.Logger, path string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Tracing: buat span untuk request ini
		ctx, span := tracer.Start(r.Context(), "HTTP "+r.Method+" "+path,
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.path", path),
			),
		)
		defer span.End()

		// Simpan context dengan trace ke request
		r = r.WithContext(ctx)

		// Catat response status code
		rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		// Jalankan handler
		next(rec, r)

		duration := time.Since(start)
		statusCode := fmt.Sprintf("%d", rec.statusCode)

		// Metrics: catat counter dan histogram
		httpRequestsTotal.WithLabelValues(r.Method, path, statusCode).Inc()
		httpDuration.WithLabelValues(r.Method, path).Observe(duration.Seconds())

		// Tracing: tambahkan status code ke span
		span.SetAttributes(attribute.Int("http.status_code", rec.statusCode))

		// Logging: structured log dengan trace_id
		traceID := span.SpanContext().TraceID().String()
		logger.Info("Request selesai",
			zap.String("method", r.Method),
			zap.String("path", path),
			zap.Int("status_code", rec.statusCode),
			zap.Duration("duration", duration),
			zap.String("trace_id", traceID),
			zap.String("client_ip", r.RemoteAddr),
		)
	}
}

// ============================================================================
// HANDLERS
// ============================================================================

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "service.GetAllAccounts")
	defer span.End()

	time.Sleep(time.Duration(20+rand.Intn(80)) * time.Millisecond)

	accounts := []map[string]interface{}{
		{"id": "acc-001", "holder": "Ayash", "balance": 500000},
		{"id": "acc-002", "holder": "Budi", "balance": 300000},
		{"id": "acc-003", "holder": "Citra", "balance": 750000},
	}

	span.SetAttributes(attribute.Int("result.count", len(accounts)))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Status: "success", Data: accounts})
}

func handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "service.GetAccountByID")
	defer span.End()

	id := r.PathValue("id")
	span.SetAttributes(attribute.String("account.id", id))

	// Simulasi query database
	_, dbSpan := tracer.Start(ctx, "repository.GetByID",
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", "SELECT"),
		),
	)
	time.Sleep(time.Duration(30+rand.Intn(150)) * time.Millisecond)
	dbSpan.End()

	account := map[string]interface{}{
		"id": id, "holder": "Ayash", "balance": 500000,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Status: "success", Data: account})
}

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "service.Transfer")
	defer span.End()

	amount := 50000 + rand.Intn(450000)
	span.SetAttributes(attribute.Int("transfer.amount", amount))

	// Simulasi validasi
	_, valSpan := tracer.Start(ctx, "service.ValidateBalance")
	time.Sleep(10 * time.Millisecond)
	valSpan.End()

	// Simulasi database transaction
	_, txSpan := tracer.Start(ctx, "repository.ExecuteTransfer",
		trace.WithAttributes(attribute.String("db.system", "postgresql")),
	)
	time.Sleep(time.Duration(50+rand.Intn(200)) * time.Millisecond)
	txSpan.End()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{
		Status:  "success",
		Message: fmt.Sprintf("Transfer Rp %d berhasil", amount),
	})
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Status: "ok", Message: "healthy"})
}

// ============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║   GRAFANA STACK — Full Observability             ║")
	fmt.Println("║   Logs (Zap→Loki) + Traces (OTel→Tempo)         ║")
	fmt.Println("║   + Metrics (Prometheus) → Grafana               ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")

	// Inisialisasi Zap Logger
	logger := initLogger()
	defer logger.Sync()

	// Inisialisasi OpenTelemetry Tracer
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Gagal inisialisasi tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())

	logger.Info("Aplikasi dimulai",
		zap.String("service", "bank-api"),
		zap.String("version", "1.0.0"),
	)

	mux := http.NewServeMux()

	// Endpoint bisnis — dibungkus observabilityMiddleware (Logs + Traces + Metrics)
	mux.HandleFunc("GET /accounts", observabilityMiddleware(logger, "/accounts", handleGetAccounts))
	mux.HandleFunc("GET /accounts/{id}", observabilityMiddleware(logger, "/accounts/{id}", handleGetAccountByID))
	mux.HandleFunc("POST /transfer", observabilityMiddleware(logger, "/transfer", handleTransfer))
	mux.HandleFunc("GET /health", observabilityMiddleware(logger, "/health", handleHealth))

	// Endpoint khusus Prometheus
	mux.Handle("GET /metrics", promhttp.Handler())

	fmt.Println("\n📊 Endpoints:")
	fmt.Println("  GET  http://localhost:8080/accounts")
	fmt.Println("  GET  http://localhost:8080/accounts/{id}")
	fmt.Println("  POST http://localhost:8080/transfer")
	fmt.Println("  GET  http://localhost:8080/health")
	fmt.Println("  GET  http://localhost:8080/metrics")
	fmt.Println("\n📈 Grafana: http://localhost:3000 (admin/admin)")
	fmt.Println("🚀 Server listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
