package dto

type BaseResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
