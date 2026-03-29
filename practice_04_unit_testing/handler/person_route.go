package handler

import (
	"net/http"
	"practice_04_unit_testing/server"
)

func (h *PersonHandler) MapRoutes() {
	h.mux.HandleFunc(server.NewAPIPath(http.MethodPost, "/persons"), h.CreatePerson)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodGet, "/persons"), h.GetAllPersons)
	h.mux.HandleFunc(server.NewAPIPath(http.MethodGet, "/persons/{id}"), h.GetPerson)
}
