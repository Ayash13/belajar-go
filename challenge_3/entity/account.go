package entity

import "time"

type Account struct {
	ID            string    `db:"id"`
	AccountHolder string    `db:"account_holder"`
	Balance       float64   `db:"balance"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}
