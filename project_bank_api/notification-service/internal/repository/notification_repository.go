package repository

import (
	"context"
	"notification-service/internal/domain"
	"notification-service/pkg/telemetry"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type NotificationRepository interface {
	Create(ctx context.Context, notif *domain.Notification) error
	UpdateStatus(ctx context.Context, id string, status string) error
	GetAll(ctx context.Context) ([]domain.Notification, error)
	GetByAccountNo(ctx context.Context, accountNo string) ([]domain.Notification, error)
}

type notificationRepository struct {
	db *sqlx.DB
}

func NewNotificationRepository(db *sqlx.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) Create(ctx context.Context, notif *domain.Notification) error {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Notification.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "notifications"),
		attribute.String("notification.event_type", notif.EventType),
		attribute.String("notification.reference_no", notif.ReferenceNo),
	)

	query := `
		INSERT INTO notifications (event_type, reference_no, account_no, payload, callback_url, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query,
		notif.EventType,
		notif.ReferenceNo,
		notif.AccountNo,
		notif.Payload,
		notif.CallbackURL,
		notif.Status,
	).Scan(&notif.ID, &notif.CreatedAt)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Notification.UpdateStatus")
	defer span.End()

	query := `UPDATE notifications SET status = $1, sent_at = NOW() WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (r *notificationRepository) GetAll(ctx context.Context) ([]domain.Notification, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Notification.GetAll")
	defer span.End()

	var notifs []domain.Notification
	err := r.db.SelectContext(ctx, &notifs, "SELECT * FROM notifications ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return notifs, nil
}

func (r *notificationRepository) GetByAccountNo(ctx context.Context, accountNo string) ([]domain.Notification, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Notification.GetByAccountNo")
	defer span.End()

	var notifs []domain.Notification
	err := r.db.SelectContext(ctx, &notifs,
		"SELECT * FROM notifications WHERE account_no = $1 ORDER BY created_at DESC", accountNo)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return notifs, nil
}
