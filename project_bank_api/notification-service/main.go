package main

import (
	"context"
	"net/http"
	"notification-service/config"
	"notification-service/internal/callback"
	"notification-service/internal/consumer"
	"notification-service/internal/handler"
	"notification-service/internal/repository"
	"notification-service/internal/service"
	"notification-service/pkg/telemetry"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	config.InitLogger()
	defer config.Log.Sync()

	tp, err := telemetry.InitTracer("notification-service")
	if err != nil {
		config.Log.Fatal("Failed to initialize tracer", zap.Error(err))
	}
	defer tp.Shutdown(context.Background())

	db := config.ConnectDB()
	defer db.Close()

	config.RunMigrations(db)

	callbackURL := os.Getenv("PARTNER_CALLBACK_URL")
	if callbackURL == "" {
		callbackURL = "http://localhost:9999/callback"
	}

	httpClient := callback.NewHTTPClient(callbackURL)
	notifRepo := repository.NewNotificationRepository(db)
	notifService := service.NewNotificationService(notifRepo, httpClient)

	// Start Kafka consumer in background
	kafkaConsumer := consumer.NewKafkaConsumer(notifService)
	go kafkaConsumer.Start(context.Background())

	// HTTP server for health + query
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"notification-service"}`))
	})

	notifHandler := handler.NewNotificationHandler(notifService)
	notifHandler.MapRoutes(mux)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8083"
	}

	config.Log.Info("Notification Service starting", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		config.Log.Fatal("Server failed to start", zap.Error(err))
	}
}
