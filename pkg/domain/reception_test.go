package domain

import (
	"testing"
	"time"
)

func TestNewReception(t *testing.T) {
	pvzID := "pvz-123"

	// Act
	reception := NewReception(pvzID)

	// Assert
	if reception == nil {
		t.Fatal("Expected reception to be created, got nil")
	}
	if reception.PVZID != pvzID {
		t.Errorf("Expected PVZID to be %s, got %s", pvzID, reception.PVZID)
	}
	if reception.Status != ReceptionStatusInProgress {
		t.Errorf("Expected status to be %s, got %s", ReceptionStatusInProgress, reception.Status)
	}
	if reception.DateTime.IsZero() {
		t.Error("Expected dateTime to be set, got zero time")
	}
	if time.Since(reception.DateTime) > time.Second {
		t.Error("Expected dateTime to be close to current time")
	}
	if reception.Products == nil {
		t.Error("Expected products to be initialized, got nil")
	}
}

func TestReception_Close(t *testing.T) {
	reception := NewReception("pvz-123")

	// Act
	err := reception.Close()

	// Assert
	if err != nil {
		t.Errorf("Expected no error when closing reception, got: %v", err)
	}
	if reception.Status != ReceptionStatusClosed {
		t.Errorf("Expected status to be %s after closing, got %s", ReceptionStatusClosed, reception.Status)
	}

	// Act again - try to close already closed reception
	err = reception.Close()

	// Assert again
	if err == nil {
		t.Error("Expected error when closing already closed reception, got nil")
	}
}

func TestReception_AddProduct(t *testing.T) {
	reception := NewReception("pvz-123")
	product, _ := NewProduct(ProductTypeElectronics, reception.ID)

	// Act
	err := reception.AddProduct(*product)

	// Assert
	if err != nil {
		t.Errorf("Expected no error when adding product to active reception, got: %v", err)
	}
	if len(reception.Products) != 1 {
		t.Errorf("Expected products length to be 1, got %d", len(reception.Products))
	}

	// Close reception
	_ = reception.Close()

	// Try to add product to closed reception
	product2, _ := NewProduct(ProductTypeClothes, reception.ID)
	err = reception.AddProduct(*product2)

	// Assert
	if err == nil {
		t.Error("Expected error when adding product to closed reception, got nil")
	}
}

func TestReception_RemoveLastProduct_LIFO(t *testing.T) {
	// Arrange
	reception := NewReception("pvz-123")

	// No products yet
	err := reception.RemoveLastProduct()
	if err == nil {
		t.Error("Expected error when removing product from empty reception, got nil")
	}

	// Add products
	product1, _ := NewProduct(ProductTypeElectronics, reception.ID)
	product2, _ := NewProduct(ProductTypeClothes, reception.ID)
	product3, _ := NewProduct(ProductTypeShoes, reception.ID)

	_ = reception.AddProduct(*product1)
	_ = reception.AddProduct(*product2)
	_ = reception.AddProduct(*product3)

	// Act - remove last product (should be product3)
	err = reception.RemoveLastProduct()

	// Assert
	if err != nil {
		t.Errorf("Expected no error when removing last product, got: %v", err)
	}
	if len(reception.Products) != 2 {
		t.Errorf("Expected products length to be 2 after removal, got %d", len(reception.Products))
	}

	// Check LIFO behavior
	if reception.Products[len(reception.Products)-1].Type != ProductTypeClothes {
		t.Errorf("Expected last product to be %s after removal, got %s",
			ProductTypeClothes, reception.Products[len(reception.Products)-1].Type)
	}

	// Close reception
	_ = reception.Close()

	// Try to remove from closed reception
	err = reception.RemoveLastProduct()

	// Assert
	if err == nil {
		t.Error("Expected error when removing product from closed reception, got nil")
	}
}
