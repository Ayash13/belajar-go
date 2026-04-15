package dto

import "time"

type AccountCreatedEvent struct {
	EventType          string    `json:"eventType"`
	AccountNo          string    `json:"accountNo"`
	CustomerID         string    `json:"customerId"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	PhoneNo            string    `json:"phoneNo"`
	PartnerReferenceNo string    `json:"partnerReferenceNo"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"createdAt"`
}

type TransferCompletedEvent struct {
	EventType            string  `json:"eventType"`
	ReferenceNo          string  `json:"referenceNo"`
	PartnerReferenceNo   string  `json:"partnerReferenceNo"`
	SourceAccountNo      string  `json:"sourceAccountNo"`
	BeneficiaryAccountNo string  `json:"beneficiaryAccountNo"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	Remark               string  `json:"remark"`
	Status               string  `json:"status"`
	TransactionDate      string  `json:"transactionDate"`
}
