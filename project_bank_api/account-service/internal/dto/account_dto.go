package dto

// ===== Account Registration =====

type DeviceInfo struct {
	OS          string `json:"os"`
	OSVersion   string `json:"osVersion"`
	Model       string `json:"model"`
	Manufacture string `json:"manufacture"`
}

type AdditionalInfo struct {
	DeviceID string `json:"deviceId"`
	Channel  string `json:"channel"`
}

type CreateAccountRequest struct {
	PartnerReferenceNo string         `json:"partnerReferenceNo"`
	CountryCode        string         `json:"countryCode"`
	CustomerID         string         `json:"customerId"`
	DeviceInfo         DeviceInfo     `json:"deviceInfo"`
	Email              string         `json:"email"`
	Lang               string         `json:"lang"`
	Locale             string         `json:"locale"`
	Name               string         `json:"name"`
	OnboardingPartner  string         `json:"onboardingPartner"`
	PhoneNo            string         `json:"phoneNo"`
	RedirectURL        string         `json:"redirectUrl"`
	Scopes             string         `json:"scopes"`
	SeamlessData       string         `json:"seamlessData"`
	SeamlessSign       string         `json:"seamlessSign"`
	State              string         `json:"state"`
	MerchantID         string         `json:"merchantId"`
	SubMerchantID      string         `json:"subMerchantId"`
	TerminalType       string         `json:"terminalType"`
	AdditionalInfo     AdditionalInfo `json:"additionalInfo"`
}

type CreateAccountResponse struct {
	ResponseCode       string         `json:"responseCode"`
	ResponseMessage    string         `json:"responseMessage"`
	ReferenceNo        string         `json:"referenceNo"`
	PartnerReferenceNo string         `json:"partnerReferenceNo"`
	AccountNo          string         `json:"accountNo"`
	AuthCode           string         `json:"authCode"`
	APIKey             string         `json:"apiKey"`
	AccountID          string         `json:"accountId"`
	State              string         `json:"state"`
	AdditionalInfo     AdditionalInfo `json:"additionalInfo"`
}
