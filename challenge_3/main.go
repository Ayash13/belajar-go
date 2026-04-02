package main

import (
	"belajar-go/challenge_3/database"
	"belajar-go/challenge_3/handler"
	"belajar-go/challenge_3/logger"
	"belajar-go/challenge_3/middleware"
	"belajar-go/challenge_3/repository"
	"belajar-go/challenge_3/service"
	"belajar-go/challenge_3/telemetry"
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	logger.Init()
	defer logger.Log.Sync()

	// Inisialisasi OpenTelemetry Tracer
	tp, err := telemetry.InitTracer()
	if err != nil {
		logger.Log.Fatal("Failed to initialize OpenTelemetry Tracer", zap.Error(err))
	}
	defer tp.Shutdown(context.Background())

	db := database.ConnectDB()
	defer db.Close()

	redisClient := database.ConnectRedis()
	defer redisClient.Close()

	accountRepo := repository.NewAccountRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	accountService := service.NewAccountService(accountRepo, transactionRepo)

	mux := http.NewServeMux()
	
	// Expose endpoint untuk Prometheus secara publik tanpa dikenakan routing bank API
	mux.Handle("/metrics", promhttp.Handler())

	accountHandler := handler.NewAccountHandler(mux, accountService)
	accountHandler.MapRoutes(redisClient)

	finalHandler := middleware.ApplyGlobalMiddlewares(redisClient, mux)

	logger.Log.Info("Server listening", zap.String("port", "8080"))
	if err := http.ListenAndServe(":8080", finalHandler); err != nil {
		logger.Log.Fatal("Server failed to start", zap.Error(err))
	}
}
