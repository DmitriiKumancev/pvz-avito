package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/application/services"
	"github.com/gorilla/mux"
)

type PVZHandler struct {
	pvzService       services.PVZService
	receptionService services.ReceptionService
}

type CreatePVZRequest struct {
	City string `json:"city"`
}

type PVZResponse struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

type PVZWithReceptionsResponse struct {
	PVZ        PVZResponse             `json:"pvz"`
	Receptions []ReceptionWithProducts `json:"receptions"`
}

type ListPVZParams struct {
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}

func NewPVZHandler(pvzService services.PVZService, receptionService services.ReceptionService) *PVZHandler {
	return &PVZHandler{
		pvzService:       pvzService,
		receptionService: receptionService,
	}
}

func (h *PVZHandler) CreatePVZ(w http.ResponseWriter, r *http.Request) {

	var req CreatePVZRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"message":"Неверный формат запроса"}`, http.StatusBadRequest)
		return
	}

	pvz, err := h.pvzService.CreatePVZ(r.Context(), req.City)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := PVZResponse{
		ID:               pvz.ID,
		RegistrationDate: pvz.RegistrationDate,
		City:             pvz.City,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (h *PVZHandler) ListPVZ(w http.ResponseWriter, r *http.Request) {
	params := extractListParams(r)

	filter := repositories.PVZFilter{
		ReceptionStartDate: params.StartDate,
		ReceptionEndDate:   params.EndDate,
		Page:               params.Page,
		Limit:              params.Limit,
	}

	pvzList, err := h.pvzService.ListPVZ(r.Context(), filter)
	if err != nil {
		http.Error(w, `{"message":"Ошибка при получении списка ПВЗ"}`, http.StatusInternalServerError)
		return
	}

	var response []PVZWithReceptionsResponse
	for _, pvz := range pvzList {
		receptions, err := h.receptionService.GetReceptionsByPVZID(r.Context(), pvz.ID)
		if err != nil {
			http.Error(w, `{"message":"Ошибка при получении данных о приемках"}`, http.StatusInternalServerError)
			return
		}

		receptionsWithProducts := make([]ReceptionWithProducts, 0, len(receptions))
		for _, reception := range receptions {
			products, err := h.receptionService.GetProductsByReceptionID(r.Context(), reception.ID)
			if err != nil {
				http.Error(w, `{"message":"Ошибка при получении данных о товарах"}`, http.StatusInternalServerError)
				return
			}

			productResponses := make([]ProductResponse, 0, len(products))
			for _, product := range products {
				productResponses = append(productResponses, ProductResponse{
					ID:          product.ID,
					DateTime:    product.DateTime,
					Type:        string(product.Type),
					ReceptionID: product.ReceptionID,
				})
			}

			receptionsWithProducts = append(receptionsWithProducts, ReceptionWithProducts{
				Reception: ReceptionResponse{
					ID:       reception.ID,
					DateTime: reception.DateTime,
					PVZID:    reception.PVZID,
					Status:   string(reception.Status),
				},
				Products: productResponses,
			})
		}

		response = append(response, PVZWithReceptionsResponse{
			PVZ: PVZResponse{
				ID:               pvz.ID,
				RegistrationDate: pvz.RegistrationDate,
				City:             pvz.City,
			},
			Receptions: receptionsWithProducts,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *PVZHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pvzID := vars["pvzId"]

	closedReception, err := h.receptionService.CloseReception(r.Context(), pvzID)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	resp := ReceptionResponse{
		ID:       closedReception.ID,
		DateTime: closedReception.DateTime,
		PVZID:    closedReception.PVZID,
		Status:   string(closedReception.Status),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *PVZHandler) DeleteLastProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pvzID := vars["pvzId"]

	err := h.receptionService.DeleteLastProduct(r.Context(), pvzID)
	if err != nil {
		http.Error(w, `{"message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"Товар успешно удален"}`))
}

func extractListParams(r *http.Request) ListPVZParams {
	query := r.URL.Query()

	var startDate, endDate *time.Time
	if startDateStr := query.Get("startDate"); startDateStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			startDate = &parsedTime
		}
	}

	if endDateStr := query.Get("endDate"); endDateStr != "" {
		if parsedTime, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			endDate = &parsedTime
		}
	}

	page := 1
	if pageStr := query.Get("page"); pageStr != "" {
		if parsedPage, err := strconv.Atoi(pageStr); err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	limit := 10
	if limitStr := query.Get("limit"); limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 30 {
			limit = parsedLimit
		}
	}

	return ListPVZParams{
		StartDate: startDate,
		EndDate:   endDate,
		Page:      page,
		Limit:     limit,
	}
}
