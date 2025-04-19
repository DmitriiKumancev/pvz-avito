package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/dkumancev/avito-pvz/pkg/application/services"
)

type ProductHandler struct {
	receptionService services.ReceptionService
}

type AddProductRequest struct {
	Type  string `json:"type"`
	PVZID string `json:"pvzId"`
}

func NewProductHandler(receptionService services.ReceptionService) *ProductHandler {
	return &ProductHandler{
		receptionService: receptionService,
	}
}

func (h *ProductHandler) AddProduct(w http.ResponseWriter, r *http.Request) {

	var req AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	product, err := h.receptionService.AddProduct(r.Context(), req.PVZID, req.Type)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := ProductResponse{
		ID:          product.ID,
		DateTime:    product.DateTime,
		Type:        string(product.Type),
		ReceptionID: product.ReceptionID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
