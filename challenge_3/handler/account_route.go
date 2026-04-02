package handler

import (
	"belajar-go/challenge_3/middleware"
	"belajar-go/challenge_3/server"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

func (h *AccountHandler) MapRoutes(rdb *redis.Client) {
	h.mux.HandleFunc(server.NewAPIPath(http.MethodPost, "/accounts"), h.CreateAccount)
	h.mux.Handle(server.NewAPIPath(http.MethodGet, "/accounts"),
		middleware.RateLimit(rdb, "account_get", 5, 5*time.Second,
			middleware.Cache(rdb, 5*time.Minute, http.HandlerFunc(h.GetAllAccounts))))

	h.mux.Handle(server.NewAPIPath(http.MethodGet, "/accounts/{id}"),
		middleware.RateLimit(rdb, "account_get", 5, 5*time.Second,
			middleware.Cache(rdb, 5*time.Minute, http.HandlerFunc(h.GetAccountByID))))
	h.mux.HandleFunc(server.NewAPIPath(http.MethodPut, "/accounts/{id}"), h.UpdateAccount)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodDelete, "/accounts/{id}"), h.DeleteAccount)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodPost, "/transfer"), h.Transfer)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodGet, "/accounts/{id}/transactions"), h.GetTransactions)
}
