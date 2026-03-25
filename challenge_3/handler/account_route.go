package handler

import (
	"belajar-go/challenge_3/server"
	"net/http"
)

func (h *AccountHandler) MapRoutes() {
	h.mux.HandleFunc(server.NewAPIPath(http.MethodPost, "/accounts"), h.CreateAccount)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodGet, "/accounts"), h.GetAllAccounts)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodGet, "/accounts/{id}"), h.GetAccountByID)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodPut, "/accounts/{id}"), h.UpdateAccount)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodDelete, "/accounts/{id}"), h.DeleteAccount)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodPost, "/transfer"), h.Transfer)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodGet, "/accounts/{id}/transactions"), h.GetTransactions)
}
