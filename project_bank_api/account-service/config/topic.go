package config

import (
	"net"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func CreateTopics() {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	conn, err := kafka.Dial("tcp", broker)
	if err != nil {
		Log.Warn("Failed to connect to Kafka for topic creation", zap.Error(err))
		return
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		Log.Warn("Failed to get Kafka controller", zap.Error(err))
		return
	}

	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		Log.Warn("Failed to connect to Kafka controller", zap.Error(err))
		return
	}
	defer controllerConn.Close()

	topics := []kafka.TopicConfig{
		{
			Topic:             "snap.account.created",
			NumPartitions:     3,
			ReplicationFactor: 1,
		},
	}

	err = controllerConn.CreateTopics(topics...)
	if err != nil {
		Log.Warn("Failed to create Kafka topics (may already exist)", zap.Error(err))
		return
	}

	Log.Info("Kafka topics created", zap.Int("count", len(topics)))
}
