package dto

import "time"

type TransferRequest struct {
	FromAccountID string  `json:"from_account_id"`
	ToAccountID   string  `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransactionResponse struct {
	ID            string    `json:"id"`
	FromAccountID string    `json:"from_account_id"`
	ToAccountID   string    `json:"to_account_id"`
	Amount        float64   `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferResponse struct {
	Message     string              `json:"message"`
	Transaction TransactionResponse `json:"transaction"`
	FromAccount AccountResponse     `json:"from_account"`
	ToAccount   AccountResponse     `json:"to_account"`
}
