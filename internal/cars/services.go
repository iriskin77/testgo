package cars

import (
	"context"

	"github.com/iriskin77/testgo/pkg/logging"
)

type ServiceCar interface {
	CreateCar(ctx context.Context, car *CarRequest) (int, error)
	UpdateCarById(ctx context.Context, carUpdate *CarUpdateRequest) (int, error)
}

type serviceCar struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo   RepositoryCar
	logger logging.Logger
}

func NewCarService(repo RepositoryCar, logger logging.Logger) *serviceCar {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceCar{
		repo:   repo,
		logger: logger}
}

func (scar *serviceCar) CreateCar(ctx context.Context, car *CarRequest) (int, error) {
	carId, err := scar.repo.CreateCar(ctx, car)

	if err != nil {
		return 0, err
	}

	return carId, nil
}

func (scar *serviceCar) UpdateCarById(ctx context.Context, carUpdate *CarUpdateRequest) (int, error) {

	carUpdatedId, err := scar.repo.UpdateCarById(ctx, carUpdate)

	if err != nil {
		return 0, err
	}

	return carUpdatedId, nil
}
