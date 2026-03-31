package mocks

import (
	"belajar-go/challenge_3/entity"
	"context"

	"github.com/stretchr/testify/mock"
)

type AccountRepository struct {
	mock.Mock
}

func (m *AccountRepository) Create(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *AccountRepository) GetAll(ctx context.Context) ([]entity.Account, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Account), args.Error(1)
}

func (m *AccountRepository) GetByID(ctx context.Context, id string) (*entity.Account, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Account), args.Error(1)
}

func (m *AccountRepository) Update(ctx context.Context, account *entity.Account) error {
	args := m.Called(ctx, account)
	return args.Error(0)
}

func (m *AccountRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
