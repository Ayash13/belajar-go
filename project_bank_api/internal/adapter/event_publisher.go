package adapter

import (
	"belajar-go/project_bank_api/config"
	"belajar-go/project_bank_api/pkg/telemetry"
	"context"
	"encoding/json"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type EventPublisher interface {
	Publish(ctx context.Context, topic string, key string, payload interface{}) error
}

type kafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(writer *kafka.Writer) EventPublisher {
	return &kafkaPublisher{writer: writer}
}

func (p *kafkaPublisher) Publish(ctx context.Context, topic string, key string, payload interface{}) error {
	ctx, span := telemetry.Tracer.Start(ctx, "kafka.Publish")
	defer span.End()

	span.SetAttributes(
		attribute.String("kafka.topic", topic),
		attribute.String("kafka.key", key),
	)

	data, err := json.Marshal(payload)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "marshal failed")
		return err
	}

	span.SetAttributes(attribute.Int("kafka.message_size", len(data)))

	msg := kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "publish failed")
		config.Log.Error("Failed to publish Kafka event",
			zap.String("topic", topic),
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	span.SetStatus(codes.Ok, "published")
	config.Log.Info("Kafka event published",
		zap.String("topic", topic),
		zap.String("key", key),
		zap.Int("size", len(data)),
	)
	return nil
}
