package dto

type TransactionHistoryRequest struct {
	AccountNo string `json:"accountNo"`
	FromDate  string `json:"fromDateTime"`
	ToDate    string `json:"toDateTime"`
}

type TransactionHistoryResponse struct {
	ResponseCode    string               `json:"responseCode"`
	ResponseMessage string               `json:"responseMessage"`
	DetailData      []TransactionDataDTO `json:"detailData"`
}

type TransactionDataDTO struct {
	ReferenceNo          string `json:"referenceNo"`
	PartnerReferenceNo   string `json:"partnerReferenceNo"`
	SourceAccountNo      string `json:"sourceAccountNo"`
	BeneficiaryAccountNo string `json:"beneficiaryAccountNo"`
	Amount               Amount `json:"amount"`
	Remark               string `json:"remark"`
	Status               string `json:"status"`
	TransactionDate      string `json:"transactionDate"`
	Type                 string `json:"type"`
}
