package cars

import (
	"context"

	"github.com/iriskin77/testgo/models"
)

type ServiceCar interface {
	CreateCar(ctx context.Context, car *models.CarRequest) (int, error)
}

type serviceCar struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo RepositoryCar
}

func NewCarService(repo RepositoryCar) *serviceCar {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceCar{repo: repo}
}

func (scar *serviceCar) CreateCar(ctx context.Context, car *models.CarRequest) (int, error) {
	carId, err := scar.repo.CreateCar(ctx, car)

	if err != nil {
		return 0, err
	}

	return carId, nil
}
