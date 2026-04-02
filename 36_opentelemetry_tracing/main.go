package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

// ============================================================================
// SETUP TRACER — Inisialisasi OpenTelemetry TracerProvider
// ============================================================================

func initTracer() (*sdktrace.TracerProvider, error) {
	// Exporter: mengirim trace data ke Tempo via OTLP HTTP
	exporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(), // tanpa TLS untuk development
	)
	if err != nil {
		return nil, err
	}

	// Resource: metadata tentang service kita
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("bank-api"),
			semconv.ServiceVersionKey.String("1.0.0"),
			attribute.String("environment", "development"),
		),
	)
	if err != nil {
		return nil, err
	}

	// TracerProvider: menyatukan exporter dan resource
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	// Set sebagai global tracer provider
	otel.SetTracerProvider(tp)
	return tp, nil
}

// ============================================================================
// SERVICE LAYER — Logika bisnis yang di-instrumentasi dengan spans
// ============================================================================

var tracer = otel.Tracer("bank-api")

type Account struct {
	ID      string `json:"id"`
	Holder  string `json:"holder"`
	Balance int    `json:"balance"`
}

// getAccountFromDB mensimulasikan query ke database dengan span terpisah
func getAccountFromDB(ctx context.Context, id string) (*Account, error) {
	// Buat child span untuk operasi database
	ctx, span := tracer.Start(ctx, "repository.GetByID",
		trace.WithAttributes(
			attribute.String("db.system", "postgresql"),
			attribute.String("db.operation", "SELECT"),
			attribute.String("db.table", "accounts"),
			attribute.String("account.id", id),
		),
	)
	defer span.End()

	// Simulasi latency database query
	time.Sleep(50 * time.Millisecond)

	// Simulasi data
	account := &Account{
		ID:      id,
		Holder:  "Ayash",
		Balance: 500000,
	}

	span.SetAttributes(attribute.Bool("cache.hit", false))
	return account, nil
}

// processGetAccount adalah layer service
func processGetAccount(ctx context.Context, id string) (*Account, error) {
	// Buat child span untuk layer service
	ctx, span := tracer.Start(ctx, "service.GetAccountByID",
		trace.WithAttributes(
			attribute.String("account.id", id),
		),
	)
	defer span.End()

	// Simulasi validasi
	time.Sleep(5 * time.Millisecond)

	// Query ke database (membuat span lagi di dalamnya)
	account, err := getAccountFromDB(ctx, id)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	return account, nil
}

// ============================================================================
// HANDLERS — HTTP handlers dengan tracing
// ============================================================================

type response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func handleGetAccountByID(w http.ResponseWriter, r *http.Request) {
	// Buat root span untuk keseluruhan request
	ctx, span := tracer.Start(r.Context(), "handler.GetAccountByID",
		trace.WithAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		),
	)
	defer span.End()

	id := r.PathValue("id")
	span.SetAttributes(attribute.String("account.id", id))

	account, err := processGetAccount(ctx, id)
	if err != nil {
		span.RecordError(err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response{Status: "error", Message: err.Error()})
		return
	}

	span.SetAttributes(attribute.String("account.holder", account.Holder))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Status: "success", Data: account})
}

func handleGetAllAccounts(w http.ResponseWriter, r *http.Request) {
	_, span := tracer.Start(r.Context(), "handler.GetAllAccounts",
		trace.WithAttributes(
			attribute.String("http.method", r.Method),
			attribute.String("http.path", r.URL.Path),
		),
	)
	defer span.End()

	time.Sleep(30 * time.Millisecond) // simulasi query

	accounts := []Account{
		{ID: "acc-001", Holder: "Ayash", Balance: 500000},
		{ID: "acc-002", Holder: "Budi", Balance: 300000},
	}

	span.SetAttributes(attribute.Int("result.count", len(accounts)))

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response{Status: "success", Data: accounts})
}

// ============================================================================

func main() {
	fmt.Println("╔══════════════════════════════════════════════════╗")
	fmt.Println("║   OPENTELEMETRY TRACING                         ║")
	fmt.Println("║   Distributed tracing untuk Go HTTP Server      ║")
	fmt.Println("╚══════════════════════════════════════════════════╝")

	// Inisialisasi tracer
	tp, err := initTracer()
	if err != nil {
		log.Fatalf("Gagal inisialisasi tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())

	fmt.Println("\n✅ OpenTelemetry Tracer berhasil diinisialisasi")
	fmt.Println("   Traces akan dikirim ke Tempo (localhost:4318)")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /accounts", handleGetAllAccounts)
	mux.HandleFunc("GET /accounts/{id}", handleGetAccountByID)

	fmt.Println("\n📍 Endpoints:")
	fmt.Println("  GET  http://localhost:8080/accounts")
	fmt.Println("  GET  http://localhost:8080/accounts/{id}")
	fmt.Println("\n🚀 Server listening on :8080")

	log.Fatal(http.ListenAndServe(":8080", mux))
}
