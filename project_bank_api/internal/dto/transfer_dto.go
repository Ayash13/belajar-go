package dto

type Amount struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type OriginatorInfo struct {
	OriginatorCustomerNo   string `json:"originatorCustomerNo"`
	OriginatorCustomerName string `json:"originatorCustomerName"`
	OriginatorBankCode     string `json:"originatorBankCode"`
}

type TransferRequest struct {
	PartnerReferenceNo   string           `json:"partnerReferenceNo"`
	Amount               Amount           `json:"amount"`
	BeneficiaryAccountNo string           `json:"beneficiaryAccountNo"`
	BeneficiaryEmail     string           `json:"beneficiaryEmail"`
	Currency             string           `json:"currency"`
	CustomerReference    string           `json:"customerReference"`
	FeeType              string           `json:"feeType"`
	Remark               string           `json:"remark"`
	SourceAccountNo      string           `json:"sourceAccountNo"`
	TransactionDate      string           `json:"transactionDate"`
	OriginatorInfos      []OriginatorInfo `json:"originatorInfos"`
	AdditionalInfo       AdditionalInfo   `json:"additionalInfo"`
}

type TransferResponse struct {
	ResponseCode         string           `json:"responseCode"`
	ResponseMessage      string           `json:"responseMessage"`
	ReferenceNo          string           `json:"referenceNo"`
	PartnerReferenceNo   string           `json:"partnerReferenceNo"`
	Amount               Amount           `json:"amount"`
	BeneficiaryAccountNo string           `json:"beneficiaryAccountNo"`
	Currency             string           `json:"currency"`
	CustomerReference    string           `json:"customerReference"`
	SourceAccount        string           `json:"sourceAccount"`
	TransactionDate      string           `json:"transactionDate"`
	OriginatorInfos      []OriginatorInfo `json:"originatorInfos"`
	AdditionalInfo       AdditionalInfo   `json:"additionalInfo"`
}
