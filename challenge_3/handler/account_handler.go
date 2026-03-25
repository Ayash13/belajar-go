package handler

import (
	"belajar-go/challenge_3/dto"
	"belajar-go/challenge_3/service"
	"encoding/json"
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
	var req dto.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: "Invalid JSON body"})
		return
	}

	if req.FromAccountID == "" || req.ToAccountID == "" || req.Amount <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: "from_account_id, to_account_id, and a positive amount are required"})
		return
	}

	resp, err := h.service.Transfer(r.Context(), req)
	if err != nil {
		if err.Error() == "insufficient balance for transfer" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: err.Error()})
			return
		}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusNotFound, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Data: resp})
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
