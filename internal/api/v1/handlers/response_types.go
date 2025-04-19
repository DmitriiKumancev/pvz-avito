package handlers

import "time"

type ReceptionResponse struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId"`
	Status   string    `json:"status"`
}

type ProductResponse struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}

type ReceptionWithProducts struct {
	Reception ReceptionResponse `json:"reception"`
	Products  []ProductResponse `json:"products"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
