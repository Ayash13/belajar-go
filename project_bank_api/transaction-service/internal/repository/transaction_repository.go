package repository

import (
	"context"
	"database/sql"
	"errors"
	"transaction-service/internal/domain"
	"transaction-service/pkg/telemetry"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TransactionRepository interface {
	Create(ctx context.Context, tx *domain.Transaction) error
	GetByReferenceNo(ctx context.Context, referenceNo string) (*domain.Transaction, error)
	GetByPartnerReferenceNo(ctx context.Context, partnerRefNo string) (*domain.Transaction, error)
	GetByAccountNo(ctx context.Context, accountNo string) ([]domain.Transaction, error)
	GetByAccountNoWithDateRange(ctx context.Context, accountNo, fromDate, toDate string) ([]domain.Transaction, error)
}

type transactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) error {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Transaction.Create")
	defer span.End()

	span.SetAttributes(
		attribute.String("db.table", "transactions"),
		attribute.String("db.operation", "INSERT"),
		attribute.String("tx.reference_no", tx.ReferenceNo),
		attribute.Float64("tx.amount", tx.Amount),
	)

	query := `
		INSERT INTO transactions (partner_reference_no, reference_no, source_account_no, beneficiary_account_no, amount, currency, remark, fee_type, transaction_date, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at`
	err := r.db.QueryRowContext(ctx, query,
		tx.PartnerReferenceNo,
		tx.ReferenceNo,
		tx.SourceAccountNo,
		tx.BeneficiaryAccountNo,
		tx.Amount,
		tx.Currency,
		tx.Remark,
		tx.FeeType,
		tx.TransactionDate,
		tx.Status,
	).Scan(&tx.ID, &tx.CreatedAt)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
	return err
}

func (r *transactionRepository) GetByReferenceNo(ctx context.Context, referenceNo string) (*domain.Transaction, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Transaction.GetByReferenceNo")
	defer span.End()

	var tx domain.Transaction
	err := r.db.GetContext(ctx, &tx, "SELECT * FROM transactions WHERE reference_no = $1", referenceNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByPartnerReferenceNo(ctx context.Context, partnerRefNo string) (*domain.Transaction, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Transaction.GetByPartnerReferenceNo")
	defer span.End()

	var tx domain.Transaction
	err := r.db.GetContext(ctx, &tx, "SELECT * FROM transactions WHERE partner_reference_no = $1", partnerRefNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) GetByAccountNo(ctx context.Context, accountNo string) ([]domain.Transaction, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Transaction.GetByAccountNo")
	defer span.End()

	var txs []domain.Transaction
	err := r.db.SelectContext(ctx, &txs,
		"SELECT * FROM transactions WHERE source_account_no = $1 OR beneficiary_account_no = $1 ORDER BY created_at DESC",
		accountNo)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Int("db.result_count", len(txs)))
	return txs, nil
}

func (r *transactionRepository) GetByAccountNoWithDateRange(ctx context.Context, accountNo, fromDate, toDate string) ([]domain.Transaction, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "repository.Transaction.GetByAccountNoWithDateRange")
	defer span.End()

	var txs []domain.Transaction
	err := r.db.SelectContext(ctx, &txs,
		`SELECT * FROM transactions 
		 WHERE (source_account_no = $1 OR beneficiary_account_no = $1) 
		 AND transaction_date >= $2::timestamp 
		 AND transaction_date <= $3::timestamp
		 ORDER BY created_at DESC`,
		accountNo, fromDate, toDate)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, err
	}

	span.SetAttributes(attribute.Int("db.result_count", len(txs)))
	return txs, nil
}
