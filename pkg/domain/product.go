package domain

import (
	"errors"
	"time"
)

// Типы товаров
const (
	ProductTypeElectronics = "электроника"
	ProductTypeClothes     = "одежда"
	ProductTypeShoes       = "обувь"
)

// допустимые типы товаров
var ValidProductTypes = map[string]bool{
	ProductTypeElectronics: true,
	ProductTypeClothes:     true,
	ProductTypeShoes:       true,
}

type Product struct {
	ID          string    `json:"id"`
	DateTime    time.Time `json:"dateTime"`
	Type        string    `json:"type"`
	ReceptionID string    `json:"receptionId"`
}

func NewProduct(productType, receptionID string) (*Product, error) {
	if !ValidProductTypes[productType] {
		return nil, errors.New("некорректный тип товара: допустимы только электроника, одежда и обувь")
	}

	return &Product{
		DateTime:    time.Now(),
		Type:        productType,
		ReceptionID: receptionID,
	}, nil
}
