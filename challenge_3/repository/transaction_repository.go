package repository

import (
	"belajar-go/challenge_3/entity"
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
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
	query := `INSERT INTO transactions (from_account_id, to_account_id, amount)
		VALUES ($1, $2, $3) RETURNING id, created_at`
	return tx.QueryRowContext(ctx, query,
		transaction.FromAccountID, transaction.ToAccountID, transaction.Amount,
	).Scan(&transaction.ID, &transaction.CreatedAt)
}

func (r *transactionRepositoryImpl) GetByAccountID(ctx context.Context, accountID string) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	query := `SELECT id, from_account_id, to_account_id, amount, created_at
		FROM transactions WHERE from_account_id = $1 OR to_account_id = $1 ORDER BY created_at`
	err := r.db.SelectContext(ctx, &transactions, query, accountID)
	return transactions, err
}

func (r *transactionRepositoryImpl) UpdateAccountBalance(ctx context.Context, tx *sqlx.Tx, accountID string, newBalance float64) error {
	query := `UPDATE accounts SET balance = $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.ExecContext(ctx, query, newBalance, accountID)
	return err
}

func (r *transactionRepositoryImpl) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return r.db.BeginTxx(ctx, nil)
}
