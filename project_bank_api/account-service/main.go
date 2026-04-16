package main

import (
	"account-service/config"
	"account-service/internal/adapter"
	"account-service/internal/grpcserver"
	"account-service/internal/handler"
	"account-service/internal/middleware"
	"account-service/internal/repository"
	"account-service/internal/service"
	"account-service/pkg/telemetry"
	"context"
	"net"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	config.InitLogger()
	defer config.Log.Sync()

	tp, err := telemetry.InitTracer("account-service")
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

	publisher := adapter.NewKafkaPublisher(kafkaWriter)

	accountRepo := repository.NewAccountRepository(db)
	accountService := service.NewAccountService(accountRepo, publisher)

	// --- gRPC Server ---
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	grpcServer := grpc.NewServer()
	grpcserver.RegisterAccountServer(grpcServer, accountService)

	go func() {
		lis, err := net.Listen("tcp", ":"+grpcPort)
		if err != nil {
			config.Log.Fatal("Failed to listen for gRPC", zap.Error(err))
		}
		config.Log.Info("gRPC server starting", zap.String("port", grpcPort))
		if err := grpcServer.Serve(lis); err != nil {
			config.Log.Fatal("gRPC server failed", zap.Error(err))
		}
	}()

	// --- HTTP Server ---
	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","service":"account-service"}`))
	})

	accountHandler := handler.NewAccountHandler(mux, accountService)
	accountHandler.MapRoutes(redisClient)

	snapMiddleware := handler.NewSNAPMiddleware(mux)
	finalHandler := middleware.ApplyGlobalMiddlewares(redisClient, snapMiddleware)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8081"
	}

	config.Log.Info("Account Service HTTP server starting", zap.String("port", port))
	if err := http.ListenAndServe(":"+port, finalHandler); err != nil {
		config.Log.Fatal("Server failed to start", zap.Error(err))
	}
}
