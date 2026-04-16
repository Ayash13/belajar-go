package handler

import (
	"encoding/json"
	"net/http"
	"notification-service/internal/service"
)

type NotificationHandler struct {
	service service.NotificationService
}

func NewNotificationHandler(service service.NotificationService) *NotificationHandler {
	return &NotificationHandler{service: service}
}

func (h *NotificationHandler) MapRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /notifications", h.GetNotifications)
}

func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accountNo := r.URL.Query().Get("accountNo")

	var notifs interface{}
	var err error

	if accountNo != "" {
		notifs, err = h.service.GetNotifications(r.Context(), accountNo)
	} else {
		notifs, err = h.service.GetAllNotifications(r.Context())
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "internal server error"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(notifs)
}
