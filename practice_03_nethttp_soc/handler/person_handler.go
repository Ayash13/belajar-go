package handler

import (
	"encoding/json"
	"net/http"
	"nethttp/dto"
	"nethttp/service"
	"strconv"
)

type PersonHandler struct {
	mux     *http.ServeMux
	service service.PersonService
}

func NewPersonHandler(mux *http.ServeMux, service service.PersonService) *PersonHandler {
	return &PersonHandler{mux: mux, service: service}
}

func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {

	var req dto.PersonCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: err.Error()})
		return
	}

	res, err := h.service.CreatePerson(r.Context(), req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusCreated, Status: "success", Data: res})
}

func (h *PersonHandler) GetAllPersons(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	persons, err := h.service.GetAllPersons(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusInternalServerError, Status: "error", Message: err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Data: persons})
}

func (h *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusBadRequest, Status: "error", Message: "Invalid ID format"})
		return
	}

	person, err := h.service.GetPerson(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusNotFound, Status: "error", Message: "Person not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.BaseResponse{Code: http.StatusOK, Status: "success", Data: person})
}
