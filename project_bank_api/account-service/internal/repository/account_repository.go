package repository

import (
	"account-service/internal/domain"
	"account-service/pkg/telemetry"
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type AccountRepository interface {
	Create(ctx context.Context, account *domain.Account) error
	GetByID(ctx context.Context, id string) (*domain.Account, error)
	GetByAccountNo(ctx context.Context, accountNo string) (*domain.Account, error)
	GetByPartnerReferenceNo(ctx context.Context, partnerRefNo string) (*domain.Account, error)
	GetAll(ctx context.Context) ([]domain.Account, error)
	UpdateBalance(ctx context.Context, accountNo string, balance float64) error
}

type accountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepository{db: db}
}

func (r *accountRepository) Create(ctx context.Context, account *domain.Account) error {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Account.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "accounts"),
		attribute.String("db.operation", "INSERT"),
		attribute.String("account.account_no", account.AccountNo),
	)

	query := `
		INSERT INTO accounts (account_no, customer_id, name, email, phone_no, balance, currency, partner_reference_no, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query,
		account.AccountNo,
		account.CustomerID,
		account.Name,
		account.Email,
		account.PhoneNo,
		account.Balance,
		account.Currency,
		account.PartnerReferenceNo,
		account.Status,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (r *accountRepository) GetByID(ctx context.Context, id string) (*domain.Account, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Account.GetByID")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "accounts"),
		attribute.String("db.operation", "SELECT"),
		attribute.String("account.id", id),
	)

	var account domain.Account
	err := r.db.GetContext(ctx, &account, "SELECT * FROM accounts WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetAttributes(attribute.Bool("db.not_found", true))
			return nil, err
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetByAccountNo(ctx context.Context, accountNo string) (*domain.Account, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Account.GetByAccountNo")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "accounts"),
		attribute.String("db.operation", "SELECT"),
		attribute.String("account.account_no", accountNo),
	)

	var account domain.Account
	err := r.db.GetContext(ctx, &account, "SELECT * FROM accounts WHERE account_no = $1", accountNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetAttributes(attribute.Bool("db.not_found", true))
			return nil, err
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetByPartnerReferenceNo(ctx context.Context, partnerRefNo string) (*domain.Account, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Account.GetByPartnerReferenceNo")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "accounts"),
		attribute.String("db.operation", "SELECT"),
		attribute.String("account.partner_reference_no", partnerRefNo),
	)

	var account domain.Account
	err := r.db.GetContext(ctx, &account, "SELECT * FROM accounts WHERE partner_reference_no = $1", partnerRefNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetAttributes(attribute.Bool("db.not_found", true))
			return nil, err
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) GetAll(ctx context.Context) ([]domain.Account, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Account.GetAll")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "accounts"),
		attribute.String("db.operation", "SELECT"),
	)

	var accounts []domain.Account
	err := r.db.SelectContext(ctx, &accounts, "SELECT * FROM accounts ORDER BY created_at DESC")
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Int("db.result_count", len(accounts)))
	return accounts, nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, accountNo string, balance float64) error {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Account.UpdateBalance")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "accounts"),
		attribute.String("db.operation", "UPDATE"),
		attribute.String("account.account_no", accountNo),
		attribute.Float64("account.new_balance", balance),
	)

	_, err := r.db.ExecContext(ctx, "UPDATE accounts SET balance = $1, updated_at = NOW() WHERE account_no = $2", balance, accountNo)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}
