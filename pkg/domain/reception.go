package domain

import (
	"errors"
	"time"
)

// Статусы приемки товаров
const (
	ReceptionStatusInProgress = "in_progress"
	ReceptionStatusClosed     = "close"
)

//  приемка товаров
type Reception struct {
	ID       string    `json:"id"`
	DateTime time.Time `json:"dateTime"`
	PVZID    string    `json:"pvzId"`
	Status   string    `json:"status"`
	Products []Product `json:"-"` // Товары, связанные с приемкой (не вклчаются в JSON напрямую)
}

func NewReception(pvzID string) *Reception {
	return &Reception{
		DateTime: time.Now(),
		PVZID:    pvzID,
		Status:   ReceptionStatusInProgress,
		Products: make([]Product, 0),
	}
}

// закрытие приемки
func (r *Reception) Close() error {
	if r.Status == ReceptionStatusClosed {
		return errors.New("приемка уже закрыта")
	}
	r.Status = ReceptionStatusClosed
	return nil
}

// check активна ли приемка
func (r *Reception) IsActive() bool {
	return r.Status == ReceptionStatusInProgress
}

func (r *Reception) AddProduct(product Product) error {
	if !r.IsActive() {
		return errors.New("нельзя добавить товар в закрытую приемку")
	}
	r.Products = append(r.Products, product)
	return nil
}

//удаление последнего добавленного товара (LIFO)
func (r *Reception) RemoveLastProduct() error {
	if !r.IsActive() {
		return errors.New("нельзя удалить товар из закрытой приемки")
	}

	if len(r.Products) == 0 {
		return errors.New("нет товаров для удаления")
	}

	r.Products = r.Products[:len(r.Products)-1]
	return nil
}
