package main

import (
	"context"
	"net/http"
	"os"
	"transaction-service/config"
	"transaction-service/internal/adapter"
	"transaction-service/internal/grpcclient"
	"transaction-service/internal/handler"
	"transaction-service/internal/middleware"
	"transaction-service/internal/repository"
	"transaction-service/internal/service"
	"transaction-service/pkg/telemetry"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	config.InitLogger()
	defer config.Log.Sync()

	tp, err := telemetry.InitTracer("transaction-service")
	if err != nil {
		config.Log.Fatal("Failed to initialize tracer", zap.Error(err))
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

	// gRPC client to account-service
	accountGrpcAddr := os.Getenv("ACCOUNT_GRPC_ADDR")
	if accountGrpcAddr == "" {
		accountGrpcAddr = "account-service:50051"
	}

	accountClient, err := grpcclient.NewAccountClient(accountGrpcAddr)
	if err != nil {
		config.Log.Fatal("Failed to connect to account-service gRPC", zap.Error(err))
	}
	defer accountClient.Close()

	publisher := adapter.NewKafkaPublisher(kafkaWriter)

	transactionRepo := repository.NewTransactionRepository(db)
	transactionService := service.NewTransactionService(transactionRepo, accountClient, publisher)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"transaction-service"}`))
	})

	transactionHandler := handler.NewTransactionHandler(mux, transactionService)
	transactionHandler.MapRoutes()

	snapMiddleware := handler.NewSNAPMiddleware(mux)
	finalHandler := middleware.ApplyGlobalMiddlewares(redisClient, snapMiddleware)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8082"
	}

	config.Log.Info("Transaction Service HTTP server starting", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		config.Log.Fatal("Server failed to start", zap.Error(err))
	}
}
