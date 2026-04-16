package consumer

import (
	"context"
	"notification-service/config"
	"notification-service/internal/service"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaConsumer struct {
	service service.NotificationService
}

func NewKafkaConsumer(service service.NotificationService) *KafkaConsumer {
	return &KafkaConsumer{service: service}
}

func (c *KafkaConsumer) Start(ctx context.Context) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	topics := []string{"snap.account.created", "snap.transfer.completed"}

	for _, topic := range topics {
		go c.consumeTopic(ctx, broker, topic)
	}
}

func (c *KafkaConsumer) consumeTopic(ctx context.Context, broker string, topic string) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		GroupID:  "notification-consumer-group",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		MaxWait:  1 * time.Second,
	})
	defer reader.Close()

	config.Log.Info("Started consuming topic", zap.String("topic", topic))

	for {
		select {
		case <-ctx.Done():
			config.Log.Info("Stopping consumer", zap.String("topic", topic))
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				config.Log.Error("Failed to read message", zap.Error(err), zap.String("topic", topic))
				time.Sleep(1 * time.Second)
				continue
			}

			config.Log.Info("Received message",
				zap.String("topic", msg.Topic),
				zap.String("key", string(msg.Key)),
			)

			err = c.service.ProcessEvent(ctx, msg.Topic, string(msg.Key), msg.Value)
			if err != nil {
				config.Log.Error("Failed to process message", zap.Error(err), zap.String("topic", msg.Topic))
			}
		}
	}
}
