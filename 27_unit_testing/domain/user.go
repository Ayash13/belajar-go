package domain

import "context"

type User struct {
	ID    int64
	Name  string
	Email string
}

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*User, error)
}

type EventPublisher interface {
	PublishEvent(ctx context.Context, event string) error
}
