package domain

import (
	"errors"
	"time"
)

// Список допустимых городов для ПВЗ
var ValidCities = map[string]bool{
	"Москва":          true,
	"Санкт-Петербург": true,
	"Казань":          true,
}

// пункт выдачи заказов
type PVZ struct {
	ID               string    `json:"id"`
	RegistrationDate time.Time `json:"registrationDate"`
	City             string    `json:"city"`
}

func NewPVZ(city string) (*PVZ, error) {
	if !ValidCities[city] {
		return nil, errors.New("город не поддерживается: разрешены только Москва, Санкт-Петербург и Казань")
	}

	return &PVZ{
		City:             city,
		RegistrationDate: time.Now(),
	}, nil
}
