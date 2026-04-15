package repository

import (
	"belajar-go/challenge_3/entity"
	"belajar-go/challenge_3/telemetry"
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"go.opentelemetry.io/otel/attribute"
)

var ErrNotFound = errors.New("not found")

type TransactionRepository interface {
	Create(ctx context.Context, tx *sqlx.Tx, transaction *entity.Transaction) error
	GetByAccountID(ctx context.Context, accountID string) ([]entity.Transaction, error)
	UpdateAccountBalance(ctx context.Context, tx *sqlx.Tx, accountID string, newBalance float64) error
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
}

type transactionRepositoryImpl struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) Create(ctx context.Context, tx *sqlx.Tx, transaction *entity.Transaction) error {
	ctx, span := telemetry.Tracer.Start(ctx, "TransactionRepository.Create")
	defer span.End()

	query := `INSERT INTO transactions (partner_reference_no, from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	
	span.SetAttributes(
		attribute.String("db.query", query),
		attribute.String("transfer.from", transaction.FromAccountID),
		attribute.String("transfer.to", transaction.ToAccountID),
	)

	return tx.QueryRowContext(ctx, query,
		transaction.PartnerReferenceNo, transaction.FromAccountID, transaction.ToAccountID, transaction.Amount,
	).Scan(&transaction.ID, &transaction.CreatedAt)
}

func (r *transactionRepositoryImpl) GetByAccountID(ctx context.Context, accountID string) ([]entity.Transaction, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "TransactionRepository.GetByAccountID")
	defer span.End()

	var transactions []entity.Transaction
	query := `SELECT id, partner_reference_no, from_account_id, to_account_id, amount, created_at
		FROM transactions WHERE from_account_id = $1 OR to_account_id = $1 ORDER BY created_at`
	
	span.SetAttributes(
		attribute.String("db.query", query),
		attribute.String("account.id", accountID),
	)

	err := r.db.SelectContext(ctx, &transactions, query, accountID)
	return transactions, err
}

func (r *transactionRepositoryImpl) UpdateAccountBalance(ctx context.Context, tx *sqlx.Tx, accountID string, newBalance float64) error {
	ctx, span := telemetry.Tracer.Start(ctx, "TransactionRepository.UpdateAccountBalance")
	defer span.End()

	query := `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`
	
	span.SetAttributes(
		attribute.String("db.query", query),
		attribute.String("account.id", accountID),
	)

	_, err := tx.ExecContext(ctx, query, newBalance, accountID)
	return err
}

func (r *transactionRepositoryImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "TransactionRepository.BeginTx")
	defer span.End()
	return r.db.BeginTxx(ctx, nil)
}
