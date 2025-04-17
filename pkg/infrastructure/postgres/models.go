package postgres

import (
	"time"

	"github.com/dkumancev/avito-pvz/pkg/domain"
)

// модель ПВЗ в БД
type PVZModel struct {
	ID               string    `db:"id"`
	RegistrationDate time.Time `db:"registration_date"`
	City             string    `db:"city"`
}

// ToEntity преобразует модель БД в доменную сущность
func (p *PVZModel) ToEntity() *domain.PVZ {
	return &domain.PVZ{
		ID:               p.ID,
		RegistrationDate: p.RegistrationDate,
		City:             p.City,
	}
}

// FromEntity преобразует доменную сущность в модель БД
func (p *PVZModel) FromEntity(pvz *domain.PVZ) {
	p.ID = pvz.ID
	p.RegistrationDate = pvz.RegistrationDate
	p.City = pvz.City
}

// модель приемки товаров в БД
type ReceptionModel struct {
	ID       string    `db:"id"`
	DateTime time.Time `db:"date_time"`
	PVZID    string    `db:"pvz_id"`
	Status   string    `db:"status"`
}

// ToEntity преобразует модель БД в доменную сущность
func (r *ReceptionModel) ToEntity() *domain.Reception {
	reception := &domain.Reception{
		ID:       r.ID,
		DateTime: r.DateTime,
		PVZID:    r.PVZID,
		Status:   r.Status,
		Products: make([]domain.Product, 0),
	}
	return reception
}

// FromEntity преобразует доменную сущность в модель БД
func (r *ReceptionModel) FromEntity(reception *domain.Reception) {
	r.ID = reception.ID
	r.DateTime = reception.DateTime
	r.PVZID = reception.PVZID
	r.Status = reception.Status
}

// модель товара в БД
type ProductModel struct {
	ID          string    `db:"id"`
	DateTime    time.Time `db:"date_time"`
	Type        string    `db:"type"`
	ReceptionID string    `db:"reception_id"`
}

// ToEntity преобразует модель БД в доменную сущность
func (p *ProductModel) ToEntity() *domain.Product {
	return &domain.Product{
		ID:          p.ID,
		DateTime:    p.DateTime,
		Type:        p.Type,
		ReceptionID: p.ReceptionID,
	}
}

// FromEntity преобразует доменную сущность в модель БД
func (p *ProductModel) FromEntity(product *domain.Product) {
	p.ID = product.ID
	p.DateTime = product.DateTime
	p.Type = product.Type
	p.ReceptionID = product.ReceptionID
}

// модель последовательности товаров в приемке
type ProductSequenceModel struct {
	ID          int    `db:"id"`
	ReceptionID string `db:"reception_id"`
	ProductID   string `db:"product_id"`
}

// модель пользователя в БД
type UserModel struct {
	ID           string    `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	Role         string    `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
}
