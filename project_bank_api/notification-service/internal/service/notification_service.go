package service

import (
	"context"
	"encoding/json"
	"notification-service/config"
	"notification-service/internal/callback"
	"notification-service/internal/domain"
	"notification-service/internal/repository"
	"notification-service/pkg/telemetry"
	"os"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type NotificationService interface {
	ProcessEvent(ctx context.Context, eventType string, key string, payload []byte) error
	GetNotifications(ctx context.Context, accountNo string) ([]domain.Notification, error)
	GetAllNotifications(ctx context.Context) ([]domain.Notification, error)
}

type notificationService struct {
	notifRepo  repository.NotificationRepository
	httpClient callback.HTTPClient
}

func NewNotificationService(
	notifRepo repository.NotificationRepository,
	httpClient callback.HTTPClient,
) NotificationService {
	return &notificationService{
		notifRepo:  notifRepo,
		httpClient: httpClient,
	}
}

func (s *notificationService) ProcessEvent(ctx context.Context, eventType string, referenceNo string, payload []byte) error {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Notification.ProcessEvent")
	defer span.End()

	span.SetAttributes(
		attribute.String("event.type", eventType),
		attribute.String("event.reference_no", referenceNo),
		attribute.Int("event.payload_size", len(payload)),
	)

	// Extract account_no from payload (we know it's in the payload struct but we simply parse it as JSON map)
	var data map[string]interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to parse event payload")
		return err
	}

	accountNo := ""
	if acc, ok := data["accountNo"].(string); ok {
		accountNo = acc
	} else if src, ok := data["sourceAccountNo"].(string); ok {
		accountNo = src // Use source account for transfer events
	}

	notif := &domain.Notification{
		EventType:   eventType,
		ReferenceNo: referenceNo,
		AccountNo:   accountNo,
		Payload:     string(payload),
		CallbackURL: os.Getenv("PARTNER_CALLBACK_URL"),
		Status:      "PENDING",
	}

	if err := s.notifRepo.Create(ctx, notif); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to save notification")
		return err
	}

	// Send callback
	statusCode, err := s.httpClient.SendCallback(ctx, payload)

	status := "SUCCESS"
	if err != nil || statusCode >= 400 {
		status = "FAILED"
	}

	if err := s.notifRepo.UpdateStatus(ctx, notif.ID, status); err != nil {
		span.RecordError(err)
		config.Log.Error("Failed to update notification status", zap.Error(err), zap.String("id", notif.ID))
	}

	if status == "FAILED" {
		span.SetStatus(codes.Error, "callback failed")
		config.Log.Warn("Callback failed", zap.String("id", notif.ID), zap.Int("status_code", statusCode))
	} else {
		span.SetStatus(codes.Ok, "callback success")
		config.Log.Info("Callback delivered", zap.String("id", notif.ID))
	}

	return nil
}

func (s *notificationService) GetNotifications(ctx context.Context, accountNo string) ([]domain.Notification, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Notification.GetNotifications")
	defer span.End()
	return s.notifRepo.GetByAccountNo(ctx, accountNo)
}

func (s *notificationService) GetAllNotifications(ctx context.Context) ([]domain.Notification, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Notification.GetAllNotifications")
	defer span.End()
	return s.notifRepo.GetAll(ctx)
}
