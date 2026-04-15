package service

import (
	"belajar-go/project_bank_api/internal/adapter"
	"belajar-go/project_bank_api/internal/constant"
	"belajar-go/project_bank_api/internal/domain"
	"belajar-go/project_bank_api/internal/dto"
	"belajar-go/project_bank_api/internal/repository"
	"belajar-go/project_bank_api/pkg"
	"belajar-go/project_bank_api/pkg/telemetry"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type AccountService interface {
	CreateAccount(ctx context.Context, req dto.CreateAccountRequest) (*dto.CreateAccountResponse, error)
	Transfer(ctx context.Context, req dto.TransferRequest) (*dto.TransferResponse, error)
	GetAllAccounts(ctx context.Context) (*dto.GetAllAccountsResponse, error)
	GetAccountByID(ctx context.Context, accountNo string) (*dto.BalanceInquiryResponse, error)
	GetTransactionHistory(ctx context.Context, req dto.TransactionHistoryRequest) (*dto.TransactionHistoryResponse, error)
}

type accountService struct {
	accountRepo     repository.AccountRepository
	transactionRepo repository.TransactionRepository
	publisher       adapter.EventPublisher
}

func NewAccountService(
	accountRepo repository.AccountRepository,
	transactionRepo repository.TransactionRepository,
	publisher adapter.EventPublisher,
) AccountService {
	return &accountService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
		publisher:       publisher,
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

	// Publish account created event
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

func (s *accountService) Transfer(ctx context.Context, req dto.TransferRequest) (*dto.TransferResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Account.Transfer")
	defer span.End()

	span.SetAttributes(
		attribute.String("snap.partner_reference_no", req.PartnerReferenceNo),
		attribute.String("snap.source_account", req.SourceAccountNo),
		attribute.String("snap.beneficiary_account", req.BeneficiaryAccountNo),
		attribute.String("snap.amount", req.Amount.Value),
		attribute.String("snap.currency", req.Currency),
	)

	if req.PartnerReferenceNo == "" || req.SourceAccountNo == "" || req.BeneficiaryAccountNo == "" || req.Amount.Value == "" {
		span.SetStatus(codes.Error, "missing mandatory fields")
		return nil, fmt.Errorf("missing mandatory fields")
	}

	existingTx, err := s.transactionRepo.GetByPartnerReferenceNo(ctx, req.PartnerReferenceNo)
	if err == nil && existingTx != nil {
		span.SetStatus(codes.Error, "duplicate partnerReferenceNo")
		return nil, fmt.Errorf("duplicate partnerReferenceNo")
	}

	sourceAccount, err := s.accountRepo.GetByAccountNo(ctx, req.SourceAccountNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "source account not found")
			return nil, fmt.Errorf("source account not found")
		}
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	beneficiaryAccount, err := s.accountRepo.GetByAccountNo(ctx, req.BeneficiaryAccountNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "destination account not found")
			return nil, fmt.Errorf("destination account not found")
		}
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	amountValue, err := strconv.ParseFloat(req.Amount.Value, 64)
	if err != nil || amountValue <= 0 {
		span.SetStatus(codes.Error, "invalid amount")
		return nil, fmt.Errorf("invalid amount")
	}

	span.SetAttributes(
		attribute.Float64("snap.parsed_amount", amountValue),
		attribute.Float64("snap.source_balance_before", sourceAccount.Balance),
		attribute.Float64("snap.beneficiary_balance_before", beneficiaryAccount.Balance),
	)

	if sourceAccount.Balance < amountValue {
		span.SetStatus(codes.Error, "insufficient funds")
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "insufficient_funds").Inc()
		return nil, fmt.Errorf("insufficient balance for transfer")
	}

	transactionDate, err := time.Parse("2006-01-02T15:04:05+07:00", req.TransactionDate)
	if err != nil {
		transactionDate = time.Now()
	}

	referenceNo := pkg.GenerateReferenceNo()

	tx := &domain.Transaction{
		PartnerReferenceNo:   req.PartnerReferenceNo,
		ReferenceNo:          referenceNo,
		SourceAccountNo:      req.SourceAccountNo,
		BeneficiaryAccountNo: req.BeneficiaryAccountNo,
		Amount:               amountValue,
		Currency:             req.Currency,
		Remark:               req.Remark,
		FeeType:              req.FeeType,
		TransactionDate:      transactionDate,
		Status:               "SUCCESS",
	}

	if err := s.transactionRepo.Create(ctx, tx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to create transaction")
		telemetry.BankTransactionsTotal.WithLabelValues("transfer", "failed").Inc()
		return nil, fmt.Errorf("internal server error")
	}

	if err := s.accountRepo.UpdateBalance(ctx, req.SourceAccountNo, sourceAccount.Balance-amountValue); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	if err := s.accountRepo.UpdateBalance(ctx, req.BeneficiaryAccountNo, beneficiaryAccount.Balance+amountValue); err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	telemetry.BankTransactionsTotal.WithLabelValues("transfer", "success").Inc()
	telemetry.BankTransactionAmountTotal.Add(amountValue)
	telemetry.BankTransactionAmountDistribution.Observe(amountValue)

	// Publish transfer completed event
	event := dto.TransferCompletedEvent{
		EventType:            "TRANSFER_COMPLETED",
		ReferenceNo:          referenceNo,
		PartnerReferenceNo:   req.PartnerReferenceNo,
		SourceAccountNo:      req.SourceAccountNo,
		BeneficiaryAccountNo: req.BeneficiaryAccountNo,
		Amount:               amountValue,
		Currency:             req.Currency,
		Remark:               req.Remark,
		Status:               "SUCCESS",
		TransactionDate:      req.TransactionDate,
	}
	_ = s.publisher.Publish(ctx, constant.TopicTransferCompleted, referenceNo, event)

	span.SetAttributes(
		attribute.Float64("snap.source_balance_after", sourceAccount.Balance-amountValue),
		attribute.Float64("snap.beneficiary_balance_after", beneficiaryAccount.Balance+amountValue),
	)
	span.SetStatus(codes.Ok, "transfer completed")

	resp := &dto.TransferResponse{
		ReferenceNo:          referenceNo,
		PartnerReferenceNo:   req.PartnerReferenceNo,
		Amount:               req.Amount,
		BeneficiaryAccountNo: req.BeneficiaryAccountNo,
		Currency:             req.Currency,
		CustomerReference:    req.CustomerReference,
		SourceAccount:        req.SourceAccountNo,
		TransactionDate:      req.TransactionDate,
		OriginatorInfos:      req.OriginatorInfos,
		AdditionalInfo:       req.AdditionalInfo,
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

	span.SetAttributes(
		attribute.String("snap.account_name", account.Name),
		attribute.Float64("snap.balance", account.Balance),
	)

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

func (s *accountService) GetTransactionHistory(ctx context.Context, req dto.TransactionHistoryRequest) (*dto.TransactionHistoryResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Account.GetTransactionHistory")
	defer span.End()

	span.SetAttributes(
		attribute.String("snap.account_no", req.AccountNo),
		attribute.String("snap.from_date", req.FromDate),
		attribute.String("snap.to_date", req.ToDate),
	)

	if req.AccountNo == "" {
		span.SetStatus(codes.Error, "missing mandatory fields")
		return nil, fmt.Errorf("missing mandatory fields")
	}

	_, err := s.accountRepo.GetByAccountNo(ctx, req.AccountNo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			span.SetStatus(codes.Error, "account not found")
			return nil, fmt.Errorf("account not found")
		}
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	var txs []domain.Transaction
	if req.FromDate != "" && req.ToDate != "" {
		txs, err = s.transactionRepo.GetByAccountNoWithDateRange(ctx, req.AccountNo, req.FromDate, req.ToDate)
	} else {
		txs, err = s.transactionRepo.GetByAccountNo(ctx, req.AccountNo)
	}
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	span.SetAttributes(attribute.Int("snap.transaction_count", len(txs)))

	var detailData []dto.TransactionDataDTO
	for _, tx := range txs {
		txType := "CREDIT"
		if tx.SourceAccountNo == req.AccountNo {
			txType = "DEBIT"
		}

		detailData = append(detailData, dto.TransactionDataDTO{
			ReferenceNo:          tx.ReferenceNo,
			PartnerReferenceNo:   tx.PartnerReferenceNo,
			SourceAccountNo:      tx.SourceAccountNo,
			BeneficiaryAccountNo: tx.BeneficiaryAccountNo,
			Amount: dto.Amount{
				Value:    fmt.Sprintf("%.2f", tx.Amount),
				Currency: tx.Currency,
			},
			Remark:          tx.Remark,
			Status:          tx.Status,
			TransactionDate: tx.TransactionDate.Format("2006-01-02T15:04:05+07:00"),
			Type:            txType,
		})
	}

	return &dto.TransactionHistoryResponse{DetailData: detailData}, nil
}
