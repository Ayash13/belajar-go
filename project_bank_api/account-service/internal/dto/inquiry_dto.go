package dto

type BalanceInquiryResponse struct {
	ResponseCode    string           `json:"responseCode"`
	ResponseMessage string           `json:"responseMessage"`
	AccountNo       string           `json:"accountNo"`
	Name            string           `json:"name"`
	AccountInfos    []AccountInfoDTO `json:"accountInfos"`
}

type AccountInfoDTO struct {
	BalanceType string         `json:"balanceType"`
	Amount      AccountBalance `json:"amount"`
	Status      string         `json:"status"`
}

type AccountBalance struct {
	Value    string `json:"value"`
	Currency string `json:"currency"`
}

type GetAllAccountsResponse struct {
	ResponseCode    string       `json:"responseCode"`
	ResponseMessage string       `json:"responseMessage"`
	Accounts        []AccountDTO `json:"accounts"`
}

type AccountDTO struct {
	AccountNo  string `json:"accountNo"`
	CustomerID string `json:"customerId"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	PhoneNo    string `json:"phoneNo"`
	Balance    string `json:"balance"`
	Currency   string `json:"currency"`
	Status     string `json:"status"`
}
