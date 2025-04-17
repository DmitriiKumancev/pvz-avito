package domain

import (
	"testing"
	"time"
)

func TestNewPVZ_ValidCity(t *testing.T) {
	// Arrange
	city := "Москва"

	// Act
	pvz, err := NewPVZ(city)

	// Assert
	if err != nil {
		t.Errorf("Expected no error for valid city, got: %v", err)
	}
	if pvz == nil {
		t.Fatal("Expected PVZ to be created, got nil")
	}
	if pvz.City != city {
		t.Errorf("Expected city to be %s, got %s", city, pvz.City)
	}
	// Проверяем, что дата регистрации установлена и не равна нулю
	if pvz.RegistrationDate.IsZero() {
		t.Error("Expected registration date to be set, got zero time")
	}
	// Проверяем, что время регистрации близко к текущему времени (с погрешностью в 1 секунду)
	if time.Since(pvz.RegistrationDate) > time.Second {
		t.Error("Expected registration date to be close to current time")
	}
}

func TestNewPVZ_InvalidCity(t *testing.T) {
	// Arrange
	invalidCities := []string{
		"Новосибирск",
		"Красноярск",
		"Екатеринбург",
		"",
	}

	// Act & Assert
	for _, city := range invalidCities {
		pvz, err := NewPVZ(city)

		if err == nil {
			t.Errorf("Expected error for invalid city %s, got nil", city)
		}
		if pvz != nil {
			t.Errorf("Expected PVZ to be nil for invalid city %s, got %+v", city, pvz)
		}
	}
}
