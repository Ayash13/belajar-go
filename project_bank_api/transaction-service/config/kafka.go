package config

import (
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func ConnectKafka() *kafka.Writer {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	w := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: 10 * time.Millisecond,
		Async:        true,
	}

	Log.Info("Kafka writer initialized", zap.String("broker", broker))
	return w
}
