package handler

import (
	"belajar-go/challenge_3/dto"
	"belajar-go/challenge_3/service"
	"encoding/json"
	"fmt"
	"net/http"
)

type AccountHandler struct {
	mux     *http.ServeMux
	service service.AccountService
}

func NewAccountHandler(mux *http.ServeMux, service service.AccountService) *AccountHandler {
	return &AccountHandler{mux: mux, service: service}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: "Invalid JSON body"})
		return
	}

	if req.AccountHolder == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: "account_holder is required"})
		return
	}

	account, err := h.service.CreateAccount(r.Context(), req)
	if err != nil {
		if err.Error() == "account with this account_number already exists" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusConflict, Status: "error", Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusCreated, Status: "success", Message: "Account created successfully", Data: account})
}

func (h *AccountHandler) GetAllAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.service.GetAllAccounts(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Data: accounts})
}

func (h *AccountHandler) GetAccountByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	account, err := h.service.GetAccountByID(r.Context(), id)
	if err != nil {
		if err.Error() == "account not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusNotFound, Status: "error", Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Data: account})
}

func (h *AccountHandler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var req dto.UpdateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: "Invalid JSON body"})
		return
	}

	account, err := h.service.UpdateAccount(r.Context(), id, req)
	if err != nil {
		if err.Error() == "account not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusNotFound, Status: "error", Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Message: "Account updated successfully", Data: account})
}

func (h *AccountHandler) DeleteAccount(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.DeleteAccount(r.Context(), id)
	if err != nil {
		if err.Error() == "account not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusNotFound, Status: "error", Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Message: "Account deleted successfully"})
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	// Set default SNAP headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-TIMESTAMP", r.Header.Get("X-TIMESTAMP"))

	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("Decode error:", err)
		w.WriteHeader(http.StatusBadRequest)
		snapErr := dto.SnapBadRequest.WithReason(err.Error()).ToResponse("17")
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	if req.SourceAccountNo == "" || req.BeneficiaryAccountNo == "" || req.Amount.Value == "" {
		w.WriteHeader(http.StatusBadRequest)
		snapErr := dto.SnapInvalidMandatoryField.WithField("source/beneficiary/amount").ToResponse("17")
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	resp, err := h.service.Transfer(r.Context(), req)
	if err != nil {
		var snapErr dto.SnapResponse
		statusCode := http.StatusInternalServerError

		if err.Error() == "insufficient balance for transfer" {
			statusCode = http.StatusForbidden
			snapErr = dto.SnapInsufficientFunds.ToResponse("17")
		} else if err.Error() == "source account not found" || err.Error() == "destination account not found" {
			statusCode = http.StatusNotFound
			snapErr = dto.SnapInvalidCardAccountCustomerVirtualAccount.ToResponse("17")
		} else { // Assume duplicate partner no if it's db constraint usually
			statusCode = http.StatusConflict
			snapErr = dto.SnapDuplicatePartnerReferenceNo.ToResponse("17")
		}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(snapErr)
		return
	}

	// Apply success code to response
	successResp := dto.SnapSuccess.ToResponse("17")
	resp.ResponseCode = successResp.ResponseCode
	resp.ResponseMessage = successResp.ResponseMessage

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *AccountHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	accountID := r.PathValue("id")

	transactions, err := h.service.GetTransactionsByAccountID(r.Context(), accountID)
	if err != nil {
		if err.Error() == "account not found" {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusNotFound, Status: "error", Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Data: transactions})
}
