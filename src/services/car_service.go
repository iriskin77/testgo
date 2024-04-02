package services

import (
	"context"

	"github.com/iriskin77/testgo/models"
	"github.com/iriskin77/testgo/src/repository"
)

type ServiceCar struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo repository.Car
}

func NewCarService(repo repository.Car) *ServiceCar {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &ServiceCar{repo: repo}
}

func (scar *ServiceCar) CreateCar(ctx context.Context, car *models.CarRequest) (int, error) {
	carId, err := scar.repo.CreateCar(ctx, car)

	if err != nil {
		return 0, err
	}

	return carId, nil
}
