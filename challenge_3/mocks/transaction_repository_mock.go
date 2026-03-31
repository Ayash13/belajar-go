package mocks

import (
	"belajar-go/challenge_3/entity"
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/mock"
)

type TransactionRepository struct {
	mock.Mock
}

func (m *TransactionRepository) Create(ctx context.Context, tx *sqlx.Tx, transaction *entity.Transaction) error {
	args := m.Called(ctx, tx, transaction)
	return args.Error(0)
}

func (m *TransactionRepository) GetByAccountID(ctx context.Context, accountID string) ([]entity.Transaction, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *TransactionRepository) UpdateAccountBalance(ctx context.Context, tx *sqlx.Tx, accountID string, newBalance float64) error {
	args := m.Called(ctx, tx, accountID, newBalance)
	return args.Error(0)
}

func (m *TransactionRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sqlx.Tx), args.Error(1)
}
