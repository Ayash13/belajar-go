package mocks

import (
	"context"

	"belajar-go/27_unit_testing/domain"
	"github.com/stretchr/testify/mock"
)

type UserRepository struct {
	mock.Mock
}

func (m *UserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

type EventPublisher struct {
	mock.Mock
}

func (m *EventPublisher) PublishEvent(ctx context.Context, event string) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
