package domain

import "time"

type Notification struct {
	ID          string     `db:"id" json:"id"`
	EventType   string     `db:"event_type" json:"eventType"`
	ReferenceNo string     `db:"reference_no" json:"referenceNo"`
	AccountNo   string     `db:"account_no" json:"accountNo"`
	Payload     string     `db:"payload" json:"payload"`
	CallbackURL string     `db:"callback_url" json:"callbackUrl"`
	Status      string     `db:"status" json:"status"`
	RetryCount  int        `db:"retry_count" json:"retryCount"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	SentAt      *time.Time `db:"sent_at" json:"sentAt"`
}
