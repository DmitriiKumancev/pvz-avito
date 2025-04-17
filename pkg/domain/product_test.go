package domain

import (
	"testing"
	"time"
)

func TestNewProduct_ValidType(t *testing.T) {
	validTypes := []string{
		ProductTypeElectronics,
		ProductTypeClothes,
		ProductTypeShoes,
	}
	receptionID := "reception-123"

	for _, productType := range validTypes {
		product, err := NewProduct(productType, receptionID)

		if err != nil {
			t.Errorf("Expected no error for valid product type %s, got: %v", productType, err)
		}
		if product == nil {
			t.Fatalf("Expected product to be created for type %s, got nil", productType)
		}
		if product.Type != productType {
			t.Errorf("Expected product type to be %s, got %s", productType, product.Type)
		}
		if product.ReceptionID != receptionID {
			t.Errorf("Expected reception ID to be %s, got %s", receptionID, product.ReceptionID)
		}
		if product.DateTime.IsZero() {
			t.Error("Expected dateTime to be set, got zero time")
		}
		if time.Since(product.DateTime) > time.Second {
			t.Error("Expected dateTime to be close to current time")
		}
	}
}

func TestNewProduct_InvalidType(t *testing.T) {
	invalidTypes := []string{
		"книги",
		"мебель",
		"продукты",
		"",
	}
	receptionID := "reception-123"

	for _, productType := range invalidTypes {
		product, err := NewProduct(productType, receptionID)

		if err == nil {
			t.Errorf("Expected error for invalid product type %s, got nil", productType)
		}
		if product != nil {
			t.Errorf("Expected product to be nil for invalid type %s, got %+v", productType, product)
		}
	}
}
