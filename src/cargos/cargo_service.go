package cargos

import (
	"context"

	"github.com/iriskin77/testgo/models"
)

type ServiceCar interface {
	CreateCar(ctx context.Context, car *models.Cargo) (int, error)
}

type serviceCargo struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo RepositoryCargo
}

func NewCargoService(repo RepositoryCargo) *serviceCargo {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceCargo{repo: repo}
}

func (cr *serviceCargo) CreateCar(ctx context.Context, car *models.Cargo) (int, error) {
	return 1, nil
}
