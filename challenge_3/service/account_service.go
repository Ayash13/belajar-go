package service

import (
	"belajar-go/challenge_3/dto"
	"belajar-go/challenge_3/entity"
	"belajar-go/challenge_3/repository"
	"context"
	"database/sql"
	"errors"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (dto.AccountResponse, error)
	GetAllAccounts(ctx context.Context) ([]dto.AccountResponse, error)
	GetAccountByID(ctx context.Context, id string) (dto.AccountResponse, error)
	UpdateAccount(ctx context.Context, id string, req dto.UpdateAccountRequest) (dto.AccountResponse, error)
	DeleteAccount(ctx context.Context, id string) error
	Transfer(ctx context.Context, req dto.TransferRequest) (dto.TransferResponse, error)
	GetTransactionsByAccountID(ctx context.Context, accountID string) ([]dto.TransactionResponse, error)
}

type accountServiceImpl struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
}

func NewAccountService(accountRepo repository.AccountRepository, transactionRepo repository.TransactionRepository) AccountService {
	return &accountServiceImpl{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

func toAccountResponse(a *entity.Account) dto.AccountResponse {
	return dto.AccountResponse{
		ID:            a.ID,
		AccountHolder: a.AccountHolder,
		Balance:       a.Balance,
		CreatedAt:     a.CreatedAt,
		UpdatedAt:     a.UpdatedAt,
	}
}

func toTransactionResponse(t *entity.Transaction) dto.TransactionResponse {
	return dto.TransactionResponse{
		ID:            t.ID,
		FromAccountID: t.FromAccountID,
		ToAccountID:   t.ToAccountID,
		Amount:        t.Amount,
		CreatedAt:     t.CreatedAt,
	}
}

func (s *accountServiceImpl) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (dto.AccountResponse, error) {
	account := &entity.Account{
		AccountHolder: req.AccountHolder,
		Balance:       req.Balance,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return dto.AccountResponse{}, err
	}

	return toAccountResponse(account), nil
}

func (s *accountServiceImpl) GetAllAccounts(ctx context.Context) ([]dto.AccountResponse, error) {
	accounts, err := s.accountRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.AccountResponse, len(accounts))
	for i := range accounts {
		responses[i] = toAccountResponse(&accounts[i])
	}
	return responses, nil
}

func (s *accountServiceImpl) GetAccountByID(ctx context.Context, id string) (dto.AccountResponse, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.AccountResponse{}, errors.New("account not found")
		}
		return dto.AccountResponse{}, err
	}
	return toAccountResponse(account), nil
}

func (s *accountServiceImpl) UpdateAccount(ctx context.Context, id string, req dto.UpdateAccountRequest) (dto.AccountResponse, error) {
	account, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.AccountResponse{}, errors.New("account not found")
		}
		return dto.AccountResponse{}, err
	}

	if req.AccountHolder != "" {
		account.AccountHolder = req.AccountHolder
	}
	if req.Balance != 0 {
		account.Balance = req.Balance
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return dto.AccountResponse{}, err
	}

	return toAccountResponse(account), nil
}

func (s *accountServiceImpl) DeleteAccount(ctx context.Context, id string) error {
	_, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("account not found")
		}
		return err
	}
	return s.accountRepo.Delete(ctx, id)
}

func (s *accountServiceImpl) Transfer(ctx context.Context, req dto.TransferRequest) (dto.TransferResponse, error) {
	fromAccount, err := s.accountRepo.GetByID(ctx, req.FromAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.TransferResponse{}, errors.New("source account not found")
		}
		return dto.TransferResponse{}, err
	}

	toAccount, err := s.accountRepo.GetByID(ctx, req.ToAccountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return dto.TransferResponse{}, errors.New("destination account not found")
		}
		return dto.TransferResponse{}, err
	}

	if fromAccount.Balance < req.Amount {
		return dto.TransferResponse{}, errors.New("insufficient balance for transfer")
	}

	tx, err := s.transactionRepo.BeginTx(ctx)
	if err != nil {
		return dto.TransferResponse{}, err
	}
	defer tx.Rollback()

	fromAccount.Balance -= req.Amount
	toAccount.Balance += req.Amount

	if err := s.transactionRepo.UpdateAccountBalance(ctx, tx, fromAccount.ID, fromAccount.Balance); err != nil {
		return dto.TransferResponse{}, err
	}
	if err := s.transactionRepo.UpdateAccountBalance(ctx, tx, toAccount.ID, toAccount.Balance); err != nil {
		return dto.TransferResponse{}, err
	}

	transaction := &entity.Transaction{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	if err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		return dto.TransferResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		return dto.TransferResponse{}, err
	}

	// Re-fetch updated accounts
	fromAccount, _ = s.accountRepo.GetByID(ctx, req.FromAccountID)
	toAccount, _ = s.accountRepo.GetByID(ctx, req.ToAccountID)

	return dto.TransferResponse{
		Message:     "Transfer successful",
		Transaction: toTransactionResponse(transaction),
		FromAccount: toAccountResponse(fromAccount),
		ToAccount:   toAccountResponse(toAccount),
	}, nil
}

func (s *accountServiceImpl) GetTransactionsByAccountID(ctx context.Context, accountID string) ([]dto.TransactionResponse, error) {
	_, err := s.accountRepo.GetByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("account not found")
		}
		return nil, err
	}

	transactions, err := s.transactionRepo.GetByAccountID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TransactionResponse, len(transactions))
	for i := range transactions {
		responses[i] = toTransactionResponse(&transactions[i])
	}
	return responses, nil
}
