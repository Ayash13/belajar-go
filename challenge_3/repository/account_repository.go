package repository

import (
	"belajar-go/challenge_3/entity"
	"context"

	"github.com/jmoiron/sqlx"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	GetAll(ctx context.Context) ([]entity.Account, error)
	GetByID(ctx context.Context, id string) (*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
	Delete(ctx context.Context, id string) error
}

type accountRepositoryImpl struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) AccountRepository {
	return &accountRepositoryImpl{db: db}
}

func (r *accountRepositoryImpl) Create(ctx context.Context, account *entity.Account) error {
	query := `INSERT INTO accounts (account_holder, balance)
		VALUES ($1, $2) RETURNING id, created_at, updated_at`
	return r.db.QueryRowContext(ctx, query,
		account.AccountHolder, account.Balance,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (r *accountRepositoryImpl) GetAll(ctx context.Context) ([]entity.Account, error) {
	var accounts []entity.Account
	query := `SELECT id, account_holder, balance, created_at, updated_at FROM accounts ORDER BY created_at`
	err := r.db.SelectContext(ctx, &accounts, query)
	return accounts, err
}

func (r *accountRepositoryImpl) GetByID(ctx context.Context, id string) (*entity.Account, error) {
	var account entity.Account
	query := `SELECT id, account_holder, balance, created_at, updated_at FROM accounts WHERE id = $1`
	err := r.db.GetContext(ctx, &account, query, id)
	return &account, err
}

func (r *accountRepositoryImpl) Update(ctx context.Context, account *entity.Account) error {
	query := `UPDATE accounts SET account_holder = $1, balance = $2, updated_at = NOW()
		WHERE id = $3 RETURNING updated_at`
	return r.db.QueryRowContext(ctx, query,
		account.AccountHolder, account.Balance, account.ID,
	).Scan(&account.UpdatedAt)
}

func (r *accountRepositoryImpl) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}


