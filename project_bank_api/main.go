package main

import (
	"belajar-go/project_bank_api/config"
	"belajar-go/project_bank_api/internal/adapter"
	"belajar-go/project_bank_api/internal/handler"
	"belajar-go/project_bank_api/internal/middleware"
	"belajar-go/project_bank_api/internal/repository"
	"belajar-go/project_bank_api/internal/service"
	"belajar-go/project_bank_api/pkg/telemetry"
	"context"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	config.InitLogger()
	defer config.Log.Sync()

	tp, err := telemetry.InitTracer()
	if err != nil {
		config.Log.Fatal("Failed to initialize OpenTelemetry Tracer", zap.Error(err))
	}
	defer tp.Shutdown(context.Background())

	db := config.ConnectDB()
	defer db.Close()

	config.RunMigrations(db)

	redisClient := config.ConnectRedis()
	defer redisClient.Close()

	kafkaWriter := config.ConnectKafka()
	defer kafkaWriter.Close()

	config.CreateTopics()

	publisher := adapter.NewKafkaPublisher(kafkaWriter)

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	accountService := service.NewAccountService(accountRepo, transactionRepo, publisher)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	accountHandler := handler.NewAccountHandler(mux, accountService)
	accountHandler.MapRoutes(redisClient)

	snapMiddleware := handler.NewSNAPMiddleware(mux)
	finalHandler := middleware.ApplyGlobalMiddlewares(redisClient, snapMiddleware)

	port := os.Getenv("SNAP_SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	config.Log.Info("SNAP Bank API server starting", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		config.Log.Fatal("Server failed to start", zap.Error(err))
	}
}
