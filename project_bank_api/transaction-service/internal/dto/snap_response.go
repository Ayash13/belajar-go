package dto

import "fmt"

type SnapResponseError struct {
	HTTPStatus      int
	CaseCode        string
	ResponseMessage string
}

type SnapResponse struct {
	ResponseCode    string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
}

func (e SnapResponseError) ToResponse(serviceCode string) SnapResponse {
	respCode := fmt.Sprintf("%03d%02s%02s", e.HTTPStatus, serviceCode, e.CaseCode)
	msg := e.ResponseMessage
	if len(msg) > 150 {
		msg = msg[:150]
	}
	return SnapResponse{
		ResponseCode:    respCode,
		ResponseMessage: msg,
	}
}

func (e SnapResponseError) WithField(field string) SnapResponseError {
	e.ResponseMessage = fmt.Sprintf("%s %s", e.ResponseMessage, field)
	return e
}

func (e SnapResponseError) WithReason(reason string) SnapResponseError {
	e.ResponseMessage = fmt.Sprintf("%s. %s", e.ResponseMessage, reason)
	return e
}

var (
	SnapSuccess                                 = SnapResponseError{HTTPStatus: 200, CaseCode: "00", ResponseMessage: "Successful"}
	SnapBadRequest                              = SnapResponseError{HTTPStatus: 400, CaseCode: "00", ResponseMessage: "Bad Request"}
	SnapInvalidFieldFormat                      = SnapResponseError{HTTPStatus: 400, CaseCode: "01", ResponseMessage: "Invalid Field Format"}
	SnapInvalidMandatoryField                   = SnapResponseError{HTTPStatus: 400, CaseCode: "02", ResponseMessage: "Invalid Mandatory Field"}
	SnapUnauthorized                            = SnapResponseError{HTTPStatus: 401, CaseCode: "00", ResponseMessage: "Unauthorized"}
	SnapInvalidTokenB2B                         = SnapResponseError{HTTPStatus: 401, CaseCode: "01", ResponseMessage: "Invalid Token (B2B)"}
	SnapInsufficientFunds                       = SnapResponseError{HTTPStatus: 403, CaseCode: "14", ResponseMessage: "Insufficient Funds"}
	SnapInvalidCardAccountCustomerVirtualAccount = SnapResponseError{HTTPStatus: 404, CaseCode: "11", ResponseMessage: "Invalid Card/Account/Customer/Virtual Account"}
	SnapConflict                                = SnapResponseError{HTTPStatus: 409, CaseCode: "00", ResponseMessage: "Conflict"}
	SnapDuplicatePartnerReferenceNo             = SnapResponseError{HTTPStatus: 409, CaseCode: "01", ResponseMessage: "Duplicate partnerReferenceNo"}
	SnapInternalServerError                     = SnapResponseError{HTTPStatus: 500, CaseCode: "01", ResponseMessage: "Internal Server Error"}
)
