// Package services содержит алиасы для совместимости со старой структурой
package services

import (
	"time"

	"github.com/dkumancev/avito-pvz/pkg/application/repositories"
	"github.com/dkumancev/avito-pvz/pkg/application/services/pvz"
	"github.com/dkumancev/avito-pvz/pkg/application/services/reception"
	"github.com/dkumancev/avito-pvz/pkg/application/services/user"
)


type (
	// PVZService интерфейс сервиса ПВЗ
	PVZService = pvz.Service

	// ReceptionService интерфейс сервиса приемок
	ReceptionService = reception.Service

	// UserService интерфейс сервиса пользователей
	UserService = user.Service
)

// Функции-конструкторы для совместимости

func NewPVZService(pvzRepo repositories.PVZRepository) PVZService {
	return pvz.New(pvzRepo)
}

func NewReceptionService(
	pvzRepo repositories.PVZRepository,
	receptionRepo repositories.ReceptionRepository,
	productRepo repositories.ProductRepository,
) ReceptionService {
	return reception.New(pvzRepo, receptionRepo, productRepo)
}

func NewUserService(
	userRepo repositories.UserRepository,
	jwtSecret []byte,
	tokenExpiry time.Duration,
) UserService {
	return user.New(userRepo, jwtSecret, tokenExpiry)
}
