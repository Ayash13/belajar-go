package dto

import "time"

type CreateAccountRequest struct {
	AccountHolder string  `json:"account_holder"`
	Balance       float64 `json:"balance"`
}

type UpdateAccountRequest struct {
	AccountHolder string  `json:"account_holder"`
	Balance       float64 `json:"balance"`
}

type AccountResponse struct {
	ID            string    `json:"id"`
	AccountHolder string    `json:"account_holder"`
	Balance       float64   `json:"balance"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
