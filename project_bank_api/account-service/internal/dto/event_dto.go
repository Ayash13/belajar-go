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
