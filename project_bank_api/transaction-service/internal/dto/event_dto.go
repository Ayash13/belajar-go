package dto

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
