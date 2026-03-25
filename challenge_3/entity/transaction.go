package entity

import "time"

type Transaction struct {
	ID            string    `db:"id"`
	FromAccountID string    `db:"from_account_id"`
	ToAccountID   string    `db:"to_account_id"`
	Amount        float64   `db:"amount"`
	CreatedAt     time.Time `db:"created_at"`
}
