package cargos

import (
	"context"

	"github.com/iriskin77/testgo/models"
)

type ServiceCar interface {
	CreateCargo(ctx context.Context, cargo *models.CargoRequest) (int, error)
}

type serviceCargo struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo RepositoryCargo
}

func NewCargoService(repo RepositoryCargo) *serviceCargo {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceCargo{repo: repo}
}

func (cr *serviceCargo) CreateCargo(ctx context.Context, cargo *models.CargoRequest) (int, error) {
	newCarId, err := cr.repo.CreateCargo(ctx, cargo)

	if err != nil {
		return 0, err
	}

	return newCarId, nil
}
