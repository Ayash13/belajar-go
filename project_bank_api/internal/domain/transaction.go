package domain

import "time"

type Transaction struct {
	ID                 string    `db:"id" json:"id"`
	PartnerReferenceNo string    `db:"partner_reference_no" json:"partnerReferenceNo"`
	ReferenceNo        string    `db:"reference_no" json:"referenceNo"`
	SourceAccountNo    string    `db:"source_account_no" json:"sourceAccountNo"`
	BeneficiaryAccountNo string  `db:"beneficiary_account_no" json:"beneficiaryAccountNo"`
	Amount             float64   `db:"amount" json:"amount"`
	Currency           string    `db:"currency" json:"currency"`
	Remark             string    `db:"remark" json:"remark"`
	FeeType            string    `db:"fee_type" json:"feeType"`
	TransactionDate    time.Time `db:"transaction_date" json:"transactionDate"`
	Status             string    `db:"status" json:"status"`
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
}
