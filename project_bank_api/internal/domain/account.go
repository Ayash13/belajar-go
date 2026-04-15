package domain

import "time"

type Account struct {
	ID                 string    `db:"id" json:"id"`
	AccountNo          string    `db:"account_no" json:"accountNo"`
	CustomerID         string    `db:"customer_id" json:"customerId"`
	Name               string    `db:"name" json:"name"`
	Email              string    `db:"email" json:"email"`
	PhoneNo            string    `db:"phone_no" json:"phoneNo"`
	Balance            float64   `db:"balance" json:"balance"`
	Currency           string    `db:"currency" json:"currency"`
	PartnerReferenceNo string    `db:"partner_reference_no" json:"partnerReferenceNo"`
	Status             string    `db:"status" json:"status"`
	CreatedAt          time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt          time.Time `db:"updated_at" json:"updatedAt"`
}
