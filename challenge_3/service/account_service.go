package service

import (
	"belajar-go/challenge_3/dto"
	"belajar-go/challenge_3/entity"
	"belajar-go/challenge_3/repository"
	"belajar-go/challenge_3/telemetry"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (dto.AccountResponse, error)
	GetAllAccounts(ctx context.Context) ([]dto.AccountResponse, error)
	GetAccountByID(ctx context.Context, id string) (dto.AccountResponse, error)
	UpdateAccount(ctx context.Context, id string, req dto.UpdateAccountRequest) (dto.AccountResponse, error)
	DeleteAccount(ctx context.Context, id string) error
	Transfer(ctx context.Context, req dto.TransferRequest) (dto.SnapTransferResponse, error)
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
		ID:                 t.ID,
		PartnerReferenceNo: t.PartnerReferenceNo,
		FromAccountID:      t.FromAccountID,
		ToAccountID:        t.ToAccountID,
		Amount:             t.Amount,
		CreatedAt:          t.CreatedAt,
	}
}

func (s *accountServiceImpl) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (dto.AccountResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "AccountService.CreateAccount")
	defer span.End()

	account := &entity.Account{
		AccountHolder: req.AccountHolder,
		Balance:       req.Balance,
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return dto.AccountResponse{}, err
	}

	telemetry.BankAccountsTotal.Inc()

	span.SetAttributes(attribute.String("account.id", account.ID))

	return toAccountResponse(account), nil
}

func (s *accountServiceImpl) GetAllAccounts(ctx context.Context) ([]dto.AccountResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "AccountService.GetAllAccounts")
	defer span.End()

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
	ctx, span := telemetry.Tracer.Start(ctx, "AccountService.GetAccountByID")
	defer span.End()
	span.SetAttributes(attribute.String("account.id", id))

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
	ctx, span := telemetry.Tracer.Start(ctx, "AccountService.UpdateAccount")
	defer span.End()
	span.SetAttributes(attribute.String("account.id", id))

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
	ctx, span := telemetry.Tracer.Start(ctx, "AccountService.DeleteAccount")
	defer span.End()
	span.SetAttributes(attribute.String("account.id", id))

	_, err := s.accountRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("account not found")
		}
		return err
	}
	return s.accountRepo.Delete(ctx, id)
}

func (s *accountServiceImpl) Transfer(ctx context.Context, req dto.TransferRequest) (dto.SnapTransferResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "AccountService.Transfer")
	defer span.End()

	var amountFloat float64
	fmt.Sscanf(req.Amount.Value, "%f", &amountFloat)

	span.SetAttributes(
		attribute.String("transfer.from", req.SourceAccountNo),
		attribute.String("transfer.to", req.BeneficiaryAccountNo),
		attribute.Float64("transfer.amount", amountFloat),
	)

	fromAccount, err := s.accountRepo.GetByID(ctx, req.SourceAccountNo)
	if err != nil {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			return dto.SnapTransferResponse{}, errors.New("source account not found")
		}
		return dto.SnapTransferResponse{}, err
	}

	toAccount, err := s.accountRepo.GetByID(ctx, req.BeneficiaryAccountNo)
	if err != nil {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		if errors.Is(err, sql.ErrNoRows) {
			return dto.SnapTransferResponse{}, errors.New("destination account not found")
		}
		return dto.SnapTransferResponse{}, err
	}

	if fromAccount.Balance < amountFloat {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		return dto.SnapTransferResponse{}, errors.New("insufficient balance for transfer")
	}

	tx, err := s.transactionRepo.BeginTx(ctx)
	if err != nil {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		return dto.SnapTransferResponse{}, err
	}
	defer tx.Rollback()

	fromAccount.Balance -= amountFloat
	toAccount.Balance += amountFloat

	if err := s.transactionRepo.UpdateAccountBalance(ctx, tx, fromAccount.ID, fromAccount.Balance); err != nil {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		return dto.SnapTransferResponse{}, err
	}
	if err := s.transactionRepo.UpdateAccountBalance(ctx, tx, toAccount.ID, toAccount.Balance); err != nil {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		return dto.SnapTransferResponse{}, err
	}

	transaction := &entity.Transaction{
		PartnerReferenceNo: req.PartnerReferenceNo,
		FromAccountID:      req.SourceAccountNo,
		ToAccountID:        req.BeneficiaryAccountNo,
		Amount:             amountFloat,
	}
	if err := s.transactionRepo.Create(ctx, tx, transaction); err != nil {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		// If duplicate partner reference no, it will throw a unique constraint error from DB
		return dto.SnapTransferResponse{}, err
	}

	if err := tx.Commit(); err != nil {
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failure").Inc()
		return dto.SnapTransferResponse{}, err
	}

	telemetry.BankTransactionsTotal.WithLabelValues("transfer", "success").Inc()
	telemetry.BankTransactionAmountTotal.Add(amountFloat)
	telemetry.BankTransactionAmountDistribution.Observe(amountFloat)

	amountValue := fmt.Sprintf("%.2f", amountFloat)

	return dto.SnapTransferResponse{
		ResponseCode:         dto.SnapSuccess.ToResponse("17").ResponseCode,
		ResponseMessage:      dto.SnapSuccess.ToResponse("17").ResponseMessage,
		ReferenceNo:          transaction.ID,
		PartnerReferenceNo:   req.PartnerReferenceNo,
		Amount:               dto.Amount{Value: amountValue, Currency: req.Amount.Currency},
		BeneficiaryAccountNo: req.BeneficiaryAccountNo,
		Currency:             req.Currency,
		CustomerReference:    req.CustomerReference,
		SourceAccount:        req.SourceAccountNo,
		TransactionDate:      req.TransactionDate,
		OriginatorInfos:      req.OriginatorInfos,
		AdditionalInfo:       req.AdditionalInfo,
	}, nil
}

func (s *accountServiceImpl) GetTransactionsByAccountID(ctx context.Context, accountID string) ([]dto.TransactionResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "AccountService.GetTransactionsByAccountID")
	defer span.End()
	span.SetAttributes(attribute.String("account.id", accountID))

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
