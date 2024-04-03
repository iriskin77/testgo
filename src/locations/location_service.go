package locations

import (
	"context"

	"github.com/iriskin77/testgo/models"
)

type ServiceLocation interface {
	CreateLocation(ctx context.Context, location *models.Location) (int, error)
	GetLocationById(ctx context.Context, id int) (*models.Location, error)
	GetLocationsList(ctx context.Context) ([]models.Location, error)
}

type serviceLocation struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo RepositoryLocation
}

func NewLocationService(repo RepositoryLocation) *serviceLocation {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceLocation{repo: repo}
}

func (sl *serviceLocation) CreateLocation(ctx context.Context, location *models.Location) (int, error) {
	newLocation, err := sl.repo.CreateLocation(ctx, location)

	if err != nil {
		return 0, err
	}

	return newLocation, nil
}

func (sl *serviceLocation) GetLocationById(ctx context.Context, id int) (*models.Location, error) {

	locationById, err := sl.repo.GetLocationById(ctx, id)

	if err != nil {
		return nil, err
	}

	return locationById, err

}

func (slc *serviceLocation) GetLocationsList(ctx context.Context) ([]models.Location, error) {
	locationsList, err := slc.repo.GetLocationsList(ctx)

	if err != nil {
		return nil, err
	}

	return locationsList, err
}
