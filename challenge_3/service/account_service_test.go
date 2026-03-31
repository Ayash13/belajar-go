package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"belajar-go/challenge_3/dto"
	"belajar-go/challenge_3/entity"
	"belajar-go/challenge_3/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAccountService_CreateAccount(t *testing.T) {
	type Mocker struct {
		accountRepo     *mocks.AccountRepository
		transactionRepo *mocks.TransactionRepository
	}

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		input     dto.CreateAccountRequest
		wantErr   bool
		expected  dto.AccountResponse
	}{
		{
			desc: "SUCCESS: Create Account",
			input: dto.CreateAccountRequest{
				AccountHolder: "Ayash",
				Balance:       1000,
			},
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.accountRepo.On(
					"Create",
					mock.Anything,
					mock.AnythingOfType("*entity.Account"),
				).Run(func(args mock.Arguments) {
					account := args.Get(1).(*entity.Account)
					account.ID = "123"
				}).Return(nil)
			},
			expected: dto.AccountResponse{
				ID:            "123",
				AccountHolder: "Ayash",
				Balance:       1000,
			},
		},
		{
			desc: "ERROR: Repository Returns Error",
			input: dto.CreateAccountRequest{
				AccountHolder: "Ayash",
				Balance:       1000,
			},
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.accountRepo.On(
					"Create",
					mock.Anything,
					mock.AnythingOfType("*entity.Account"),
				).Return(assert.AnError)
			},
			expected: dto.AccountResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			accountRepo := new(mocks.AccountRepository)
			transactionRepo := new(mocks.TransactionRepository)
			mocker := &Mocker{accountRepo: accountRepo, transactionRepo: transactionRepo}

			tc.mockSetup(mocker)

			service := NewAccountService(accountRepo, transactionRepo)

			actual, err := service.CreateAccount(context.Background(), tc.input)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			}

			accountRepo.AssertExpectations(t)
			transactionRepo.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAccountByID(t *testing.T) {
	type Mocker struct {
		accountRepo     *mocks.AccountRepository
		transactionRepo *mocks.TransactionRepository
	}

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		inputID   string
		wantErr   bool
		expected  dto.AccountResponse
	}{
		{
			desc:    "SUCCESS: Get Account by ID",
			inputID: "123",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.accountRepo.On(
					"GetByID",
					mock.Anything,
					"123",
				).Return(&entity.Account{
					ID:            "123",
					AccountHolder: "Ayash",
					Balance:       1000,
				}, nil)
			},
			expected: dto.AccountResponse{
				ID:            "123",
				AccountHolder: "Ayash",
				Balance:       1000,
			},
		},
		{
			desc:    "ERROR: Account Not Found",
			inputID: "999",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.accountRepo.On(
					"GetByID",
					mock.Anything,
					"999",
				).Return(nil, sql.ErrNoRows)
			},
			expected: dto.AccountResponse{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			accountRepo := new(mocks.AccountRepository)
			transactionRepo := new(mocks.TransactionRepository)
			mocker := &Mocker{accountRepo: accountRepo, transactionRepo: transactionRepo}

			tc.mockSetup(mocker)

			service := NewAccountService(accountRepo, transactionRepo)

			actual, err := service.GetAccountByID(context.Background(), tc.inputID)

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			}

			accountRepo.AssertExpectations(t)
			transactionRepo.AssertExpectations(t)
		})
	}
}

func TestAccountService_GetAllAccounts(t *testing.T) {
	type Mocker struct {
		accountRepo     *mocks.AccountRepository
		transactionRepo *mocks.TransactionRepository
	}

	testCases := []struct {
		desc      string
		mockSetup func(m *Mocker)
		wantErr   bool
		expected  []dto.AccountResponse
	}{
		{
			desc:    "SUCCESS: Get All Accounts",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.accountRepo.On(
					"GetAll",
					mock.Anything,
				).Return([]entity.Account{
					{ID: "123", AccountHolder: "Ayash", Balance: 1000},
					{ID: "456", AccountHolder: "John", Balance: 2000},
				}, nil)
			},
			expected: []dto.AccountResponse{
				{ID: "123", AccountHolder: "Ayash", Balance: 1000},
				{ID: "456", AccountHolder: "John", Balance: 2000},
			},
		},
		{
			desc:    "SUCCESS: Get All Accounts Empty List",
			wantErr: false,
			mockSetup: func(m *Mocker) {
				m.accountRepo.On(
					"GetAll",
					mock.Anything,
				).Return([]entity.Account{}, nil)
			},
			expected: []dto.AccountResponse{},
		},
		{
			desc:    "ERROR: Repository Returns Error",
			wantErr: true,
			mockSetup: func(m *Mocker) {
				m.accountRepo.On(
					"GetAll",
					mock.Anything,
				).Return(nil, errors.New("internal error"))
			},
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			accountRepo := new(mocks.AccountRepository)
			transactionRepo := new(mocks.TransactionRepository)
			mocker := &Mocker{accountRepo: accountRepo, transactionRepo: transactionRepo}

			tc.mockSetup(mocker)

			service := NewAccountService(accountRepo, transactionRepo)

			actual, err := service.GetAllAccounts(context.Background())

			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expected, actual)
			}

			accountRepo.AssertExpectations(t)
			transactionRepo.AssertExpectations(t)
		})
	}
}
