package services

import (
	"context"

	"github.com/iriskin77/testgo/models"
	"github.com/iriskin77/testgo/src/repository"
)

type ServiceLocation struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo repository.Location
}

func NewLocationService(repo repository.Location) *ServiceLocation {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &ServiceLocation{repo: repo}
}

func (sl *ServiceLocation) CreateLocation(ctx context.Context, location *models.Location) (int, error) {
	newLocation, err := sl.repo.CreateLocation(ctx, location)

	if err != nil {
		return 0, err
	}

	return newLocation, nil
}

func (sl *ServiceLocation) GetLocationById(ctx context.Context, id int) (*models.Location, error) {

	locationById, err := sl.repo.GetLocationById(ctx, id)

	if err != nil {
		return nil, err
	}

	return locationById, err

}

func (slc *ServiceLocation) GetLocationsList(ctx context.Context) ([]models.Location, error) {
	return nil, nil
}
