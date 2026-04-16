package service

import (
	"account-service/internal/adapter"
	"account-service/internal/constant"
	"account-service/internal/domain"
	"account-service/internal/dto"
	"account-service/internal/repository"
	"account-service/pkg"
	"account-service/pkg/telemetry"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
	GetAllAccounts(ctx context.Context) (*dto.GetAllAccountsResponse, error)
	GetAccountByID(ctx context.Context, accountNo string) (*dto.BalanceInquiryResponse, error)
	// gRPC methods
	GetBalance(ctx context.Context, accountNo string) (*domain.Account, error)
	UpdateBalance(ctx context.Context, accountNo string, newBalance float64) error
}

type accountService struct {
	accountRepo repository.AccountRepository
	publisher   adapter.EventPublisher
}

func NewAccountService(
	accountRepo repository.AccountRepository,
	publisher adapter.EventPublisher,
) AccountService {
	return &accountService{
		accountRepo: accountRepo,
		publisher:   publisher,
	}
}

func (s *accountService) CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.CreateAccountResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Account.CreateAccount")
	defer span.End()

	span.SetAttributes(
		attribute.String("snap.partner_reference_no", req.PartnerReferenceNo),
		attribute.String("snap.customer_id", req.CustomerID),
		attribute.String("snap.name", req.Name),
	)

	if req.PartnerReferenceNo == "" || req.Name == "" || req.Email == "" {
		span.SetStatus(codes.Error, "missing mandatory fields")
		return nil, fmt.Errorf("missing mandatory fields")
	}

	existing, err := s.accountRepo.GetByPartnerReferenceNo(ctx, req.PartnerReferenceNo)
	if err == nil && existing != nil {
		span.SetStatus(codes.Error, "duplicate partnerReferenceNo")
		return nil, fmt.Errorf("duplicate partnerReferenceNo")
	}

	accountNo := fmt.Sprintf("88880%010d", time.Now().UnixNano()%10000000000)
	accountID := pkg.GenerateAccountID()

	span.SetAttributes(attribute.String("snap.generated_account_no", accountNo))

	account := &domain.Account{
		AccountNo:          accountNo,
		CustomerID:         req.CustomerID,
		Name:               req.Name,
		Email:              req.Email,
		PhoneNo:            req.PhoneNo,
		Balance:            0,
		Currency:           "IDR",
		PartnerReferenceNo: req.PartnerReferenceNo,
		Status:             "ACTIVE",
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create account")
		telemetry.BankTransactionsTotal.WithLabelValues("account_creation", "failed").Inc()
		return nil, fmt.Errorf("internal server error")
	}

	telemetry.BankTransactionsTotal.WithLabelValues("account_creation", "success").Inc()
	telemetry.BankAccountsTotal.Inc()

	event := dto.AccountCreatedEvent{
		EventType:          "ACCOUNT_CREATED",
		AccountNo:          accountNo,
		CustomerID:         req.CustomerID,
		Name:               req.Name,
		Email:              req.Email,
		PhoneNo:            req.PhoneNo,
		PartnerReferenceNo: req.PartnerReferenceNo,
		Status:             "ACTIVE",
		CreatedAt:          time.Now(),
	}
	_ = s.publisher.Publish(ctx, constant.TopicAccountCreated, accountNo, event)

	span.SetStatus(codes.Ok, "account created successfully")

	resp := &dto.CreateAccountResponse{
		ReferenceNo:        pkg.GenerateReferenceNo(),
		PartnerReferenceNo: req.PartnerReferenceNo,
		AccountNo:          accountNo,
		AuthCode:           pkg.GenerateAuthCode(),
		APIKey:             pkg.GenerateAPIKey(),
		AccountID:          accountID,
		State:              req.State,
		AdditionalInfo:     req.AdditionalInfo,
	}

	return resp, nil
}

func (s *accountService) GetAllAccounts(ctx context.Context) (*dto.GetAllAccountsResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Account.GetAllAccounts")
	defer span.End()

	accounts, err := s.accountRepo.GetAll(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to get accounts")
		return nil, fmt.Errorf("internal server error")
	}

	span.SetAttributes(attribute.Int("snap.account_count", len(accounts)))

	var accountDTOs []dto.AccountDTO
	for _, a := range accounts {
		accountDTOs = append(accountDTOs, dto.AccountDTO{
			AccountNo:  a.AccountNo,
			CustomerID: a.CustomerID,
			Name:       a.Name,
			Email:      a.Email,
			PhoneNo:    a.PhoneNo,
			Balance:    fmt.Sprintf("%.2f", a.Balance),
			Currency:   a.Currency,
			Status:     a.Status,
		})
	}

	return &dto.GetAllAccountsResponse{Accounts: accountDTOs}, nil
}

func (s *accountService) GetAccountByID(ctx context.Context, accountNo string) (*dto.BalanceInquiryResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Account.GetAccountByID")
	defer span.End()

	span.SetAttributes(attribute.String("snap.account_no", accountNo))

	account, err := s.accountRepo.GetByAccountNo(ctx, accountNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "account not found")
			return nil, fmt.Errorf("account not found")
		}
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	resp := &dto.BalanceInquiryResponse{
		AccountNo: account.AccountNo,
		Name:      account.Name,
		AccountInfos: []dto.AccountInfoDTO{
			{
				BalanceType: "AVAILABLE",
				Amount: dto.AccountBalance{
					Value:    fmt.Sprintf("%.2f", account.Balance),
					Currency: account.Currency,
				},
				Status: account.Status,
			},
		},
	}

	return resp, nil
}

// GetBalance is used by gRPC server to serve balance queries from transaction-service
func (s *accountService) GetBalance(ctx context.Context, accountNo string) (*domain.Account, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Account.GetBalance.gRPC")
	defer span.End()

	account, err := s.accountRepo.GetByAccountNo(ctx, accountNo)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	return account, nil
}

// UpdateBalance is used by gRPC server to update balance from transaction-service
func (s *accountService) UpdateBalance(ctx context.Context, accountNo string, newBalance float64) error {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Account.UpdateBalance.gRPC")
	defer span.End()

	return s.accountRepo.UpdateBalance(ctx, accountNo, newBalance)
}
