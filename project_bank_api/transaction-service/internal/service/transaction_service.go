package service

import (
	"context"
	"fmt"
	"strconv"
	"time"
	"transaction-service/internal/adapter"
	"transaction-service/internal/constant"
	"transaction-service/internal/domain"
	"transaction-service/internal/dto"
	"transaction-service/internal/grpcclient"
	"transaction-service/internal/repository"
	"transaction-service/pkg"
	"transaction-service/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TransactionService interface {
	Transfer(ctx context.Context, req dto.TransferRequest) (*dto.TransferResponse, error)
	GetTransactionHistory(ctx context.Context, req dto.TransactionHistoryRequest) (*dto.TransactionHistoryResponse, error)
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	accountClient   *grpcclient.AccountClient
	publisher       adapter.EventPublisher
}

func NewTransactionService(
	transactionRepo repository.TransactionRepository,
	accountClient *grpcclient.AccountClient,
	publisher adapter.EventPublisher,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		accountClient:   accountClient,
		publisher:       publisher,
	}
}

func (s *transactionService) Transfer(ctx context.Context, req dto.TransferRequest) (*dto.TransferResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Transaction.Transfer")
	defer span.End()

	span.SetAttributes(
		attribute.String("snap.partner_reference_no", req.PartnerReferenceNo),
		attribute.String("snap.source_account", req.SourceAccountNo),
		attribute.String("snap.beneficiary_account", req.BeneficiaryAccountNo),
		attribute.String("snap.amount", req.Amount.Value),
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

	// Get source balance via gRPC
	sourceResp, err := s.accountClient.GetBalance(ctx, req.SourceAccountNo)
	if err != nil || !sourceResp.Found {
		span.SetStatus(codes.Error, "source account not found")
		return nil, fmt.Errorf("source account not found")
	}

	// Get beneficiary balance via gRPC
	beneficiaryResp, err := s.accountClient.GetBalance(ctx, req.BeneficiaryAccountNo)
	if err != nil || !beneficiaryResp.Found {
		span.SetStatus(codes.Error, "destination account not found")
		return nil, fmt.Errorf("destination account not found")
	}

	amountValue, err := strconv.ParseFloat(req.Amount.Value, 64)
	if err != nil || amountValue <= 0 {
		span.SetStatus(codes.Error, "invalid amount")
		return nil, fmt.Errorf("invalid amount")
	}

	span.SetAttributes(
		attribute.Float64("snap.parsed_amount", amountValue),
		attribute.Float64("snap.source_balance_before", sourceResp.Balance),
		attribute.Float64("snap.beneficiary_balance_before", beneficiaryResp.Balance),
	)

	if sourceResp.Balance < amountValue {
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

	// Update balances via gRPC
	_, err = s.accountClient.UpdateBalance(ctx, req.SourceAccountNo, sourceResp.Balance-amountValue)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	_, err = s.accountClient.UpdateBalance(ctx, req.BeneficiaryAccountNo, beneficiaryResp.Balance+amountValue)
	if err != nil {
		span.RecordError(err)
		return nil, fmt.Errorf("internal server error")
	}

	telemetry.BankTransactionsTotal.WithLabelValues("transfer", "success").Inc()
	telemetry.BankTransactionAmountTotal.Add(amountValue)
	telemetry.BankTransactionAmountDistribution.Observe(amountValue)

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

func (s *transactionService) GetTransactionHistory(ctx context.Context, req dto.TransactionHistoryRequest) (*dto.TransactionHistoryResponse, error) {
	ctx, span := telemetry.Tracer.Start(ctx, "service.Transaction.GetTransactionHistory")
	defer span.End()

	if req.AccountNo == "" {
		span.SetStatus(codes.Error, "missing mandatory fields")
		return nil, fmt.Errorf("missing mandatory fields")
	}

	// Verify account exists via gRPC
	accountResp, err := s.accountClient.GetBalance(ctx, req.AccountNo)
	if err != nil || !accountResp.Found {
		span.SetStatus(codes.Error, "account not found")
		return nil, fmt.Errorf("account not found")
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
