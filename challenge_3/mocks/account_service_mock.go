package mocks

import (
	"context"

	"belajar-go/challenge_3/dto"

	"github.com/stretchr/testify/mock"
)

type AccountService struct {
	mock.Mock
}

func (m *AccountService) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (dto.AccountResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.AccountResponse), args.Error(1)
}

func (m *AccountService) GetAllAccounts(ctx context.Context) ([]dto.AccountResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.AccountResponse), args.Error(1)
}

func (m *AccountService) GetAccountByID(ctx context.Context, id string) (dto.AccountResponse, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(dto.AccountResponse), args.Error(1)
}

func (m *AccountService) UpdateAccount(ctx context.Context, id string, req dto.UpdateAccountRequest) (dto.AccountResponse, error) {
	args := m.Called(ctx, id, req)
	return args.Get(0).(dto.AccountResponse), args.Error(1)
}

func (m *AccountService) DeleteAccount(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *AccountService) Transfer(ctx context.Context, req dto.TransferRequest) (dto.TransferResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(dto.TransferResponse), args.Error(1)
}

func (m *AccountService) GetTransactionsByAccountID(ctx context.Context, accountID string) ([]dto.TransactionResponse, error) {
	args := m.Called(ctx, accountID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.TransactionResponse), args.Error(1)
}
