package handler

import (
	"encoding/json"
	"net/http"
	"transaction-service/internal/constant"
	"transaction-service/internal/dto"
	"transaction-service/internal/service"
	"transaction-service/pkg/telemetry"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type TransactionHandler struct {
	mux     *http.ServeMux
	service service.TransactionService
}

func NewTransactionHandler(mux *http.ServeMux, service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{mux: mux, service: service}
}

func (h *TransactionHandler) MapRoutes() {
	h.mux.HandleFunc("POST /snap/v1/transfer-intrabank", h.Transfer)
	h.mux.HandleFunc("POST /snap/v1/transaction-history", h.GetTransactionHistory)
}

func (h *TransactionHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.Tracer.Start(r.Context(), "handler.Transfer")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-TIMESTAMP", r.Header.Get("X-TIMESTAMP"))

	span.SetAttributes(
		attribute.String("snap.endpoint", "/snap/v1/transfer-intrabank"),
		attribute.String("snap.partner_id", r.Header.Get("X-PARTNER-ID")),
	)

	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		snapErr := dto.SnapBadRequest.WithReason(err.Error()).ToResponse(constant.ServiceCodeTransfer)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if req.SourceAccountNo == "" || req.BeneficiaryAccountNo == "" || req.Amount.Value == "" {
		span.SetStatus(codes.Error, "missing mandatory fields")
		w.WriteHeader(http.StatusBadRequest)
		snapErr := dto.SnapInvalidMandatoryField.WithField("source/beneficiary/amount").ToResponse(constant.ServiceCodeTransfer)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	span.SetAttributes(
		attribute.String("snap.source_account", req.SourceAccountNo),
		attribute.String("snap.beneficiary_account", req.BeneficiaryAccountNo),
		attribute.String("snap.amount", req.Amount.Value),
	)

	resp, err := h.service.Transfer(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		snapErr, statusCode := mapErrorToSnap(err, constant.ServiceCodeTransfer)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	successResp := dto.SnapSuccess.ToResponse(constant.ServiceCodeTransfer)
	resp.ResponseCode = successResp.ResponseCode
	resp.ResponseMessage = successResp.ResponseMessage

	span.SetStatus(codes.Ok, "transfer completed")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *TransactionHandler) GetTransactionHistory(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.Tracer.Start(r.Context(), "handler.GetTransactionHistory")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-TIMESTAMP", r.Header.Get("X-TIMESTAMP"))

	var req dto.TransactionHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		snapErr := dto.SnapBadRequest.WithReason(err.Error()).ToResponse(constant.ServiceCodeTransactionHistory)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	resp, err := h.service.GetTransactionHistory(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		snapErr, statusCode := mapErrorToSnap(err, constant.ServiceCodeTransactionHistory)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	successResp := dto.SnapSuccess.ToResponse(constant.ServiceCodeTransactionHistory)
	resp.ResponseCode = successResp.ResponseCode
	resp.ResponseMessage = successResp.ResponseMessage

	span.SetStatus(codes.Ok, "history retrieved")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func mapErrorToSnap(err error, serviceCode string) (dto.SnapResponse, int) {
	switch err.Error() {
	case "missing mandatory fields":
		return dto.SnapInvalidMandatoryField.ToResponse(serviceCode), http.StatusBadRequest
	case "invalid amount":
		return dto.SnapInvalidFieldFormat.WithField("amount").ToResponse(serviceCode), http.StatusBadRequest
	case "insufficient balance for transfer":
		return dto.SnapInsufficientFunds.ToResponse(serviceCode), http.StatusForbidden
	case "source account not found", "destination account not found", "account not found":
		return dto.SnapInvalidCardAccountCustomerVirtualAccount.ToResponse(serviceCode), http.StatusNotFound
	case "duplicate partnerReferenceNo":
		return dto.SnapDuplicatePartnerReferenceNo.ToResponse(serviceCode), http.StatusConflict
	default:
		return dto.SnapInternalServerError.ToResponse(serviceCode), http.StatusInternalServerError
	}
}
