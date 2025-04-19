package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dkumancev/avito-pvz/pkg/application/services"
)

type ReceptionHandler struct {
	receptionService services.ReceptionService
}

type CreateReceptionRequest struct {
	PVZID string `json:"pvzId"`
}

func NewReceptionHandler(receptionService services.ReceptionService) *ReceptionHandler {
	return &ReceptionHandler{
		receptionService: receptionService,
	}
}

func (h *ReceptionHandler) CreateReception(w http.ResponseWriter, r *http.Request) {

	var req CreateReceptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	reception, err := h.receptionService.CreateReception(r.Context(), req.PVZID)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := ReceptionResponse{
		ID:       reception.ID,
		DateTime: reception.DateTime,
		PVZID:    reception.PVZID,
		Status:   string(reception.Status),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
