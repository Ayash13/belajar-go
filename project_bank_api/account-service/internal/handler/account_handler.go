package handler

import (
	"account-service/internal/constant"
	"account-service/internal/dto"
	"account-service/internal/middleware"
	"account-service/internal/service"
	"account-service/pkg/telemetry"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type AccountHandler struct {
	mux     *http.ServeMux
	service service.AccountService
}

func NewAccountHandler(mux *http.ServeMux, service service.AccountService) *AccountHandler {
	return &AccountHandler{mux: mux, service: service}
}

func (h *AccountHandler) MapRoutes(rdb *redis.Client) {
	h.mux.HandleFunc("POST /snap/v1/account-creation", h.CreateAccount)
	h.mux.Handle("GET /snap/v1/accounts", middleware.Cache(rdb, 30*time.Second, http.HandlerFunc(h.GetAllAccounts)))
	h.mux.HandleFunc("GET /snap/v1/accounts/{accountNo}", h.GetAccountByID)
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.Tracer.Start(r.Context(), "handler.CreateAccount")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-TIMESTAMP", r.Header.Get("X-TIMESTAMP"))

	span.SetAttributes(
		attribute.String("snap.endpoint", "/snap/v1/account-creation"),
		attribute.String("snap.partner_id", r.Header.Get("X-PARTNER-ID")),
		attribute.String("snap.external_id", r.Header.Get("X-EXTERNAL-ID")),
	)

	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "invalid request body")
		w.WriteHeader(http.StatusBadRequest)
		snapErr := dto.SnapBadRequest.WithReason(err.Error()).ToResponse(constant.ServiceCodeAccountCreation)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	span.SetAttributes(attribute.String("snap.partner_reference_no", req.PartnerReferenceNo))

	resp, err := h.service.CreateAccount(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		snapErr, statusCode := mapErrorToSnap(err, constant.ServiceCodeAccountCreation)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	successResp := dto.SnapSuccess.ToResponse(constant.ServiceCodeAccountCreation)
	resp.ResponseCode = successResp.ResponseCode
	resp.ResponseMessage = successResp.ResponseMessage

	span.SetStatus(codes.Ok, "account created")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.Tracer.Start(r.Context(), "handler.GetAllAccounts")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-TIMESTAMP", r.Header.Get("X-TIMESTAMP"))

	resp, err := h.service.GetAllAccounts(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		snapErr, statusCode := mapErrorToSnap(err, constant.ServiceCodeBalanceInquiry)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	successResp := dto.SnapSuccess.ToResponse(constant.ServiceCodeBalanceInquiry)
	resp.ResponseCode = successResp.ResponseCode
	resp.ResponseMessage = successResp.ResponseMessage

	span.SetStatus(codes.Ok, "accounts retrieved")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.Tracer.Start(r.Context(), "handler.GetAccountByID")
	defer span.End()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-TIMESTAMP", r.Header.Get("X-TIMESTAMP"))

	accountNo := r.PathValue("accountNo")
	if accountNo == "" {
		span.SetStatus(codes.Error, "missing accountNo")
		w.WriteHeader(http.StatusBadRequest)
		snapErr := dto.SnapInvalidMandatoryField.WithField("accountNo").ToResponse(constant.ServiceCodeBalanceInquiry)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	span.SetAttributes(attribute.String("snap.account_no", accountNo))

	resp, err := h.service.GetAccountByID(ctx, accountNo)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		snapErr, statusCode := mapErrorToSnap(err, constant.ServiceCodeBalanceInquiry)
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	successResp := dto.SnapSuccess.ToResponse(constant.ServiceCodeBalanceInquiry)
	resp.ResponseCode = successResp.ResponseCode
	resp.ResponseMessage = successResp.ResponseMessage

	span.SetStatus(codes.Ok, "account retrieved")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func mapErrorToSnap(err error, serviceCode string) (dto.SnapResponse, int) {
	switch err.Error() {
	case "missing mandatory fields":
		return dto.SnapInvalidMandatoryField.ToResponse(serviceCode), http.StatusBadRequest
	case "account not found":
		return dto.SnapInvalidCardAccountCustomerVirtualAccount.ToResponse(serviceCode), http.StatusNotFound
	case "duplicate partnerReferenceNo":
		return dto.SnapDuplicatePartnerReferenceNo.ToResponse(serviceCode), http.StatusConflict
	default:
		return dto.SnapInternalServerError.ToResponse(serviceCode), http.StatusInternalServerError
	}
}
