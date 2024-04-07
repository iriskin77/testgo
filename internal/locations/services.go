package locations

import (
	"context"

	"github.com/iriskin77/testgo/internal/middleware"
	"github.com/iriskin77/testgo/pkg/logging"
)

type ServiceLocation interface {
	CreateLocation(ctx context.Context, location *Location) (int, error)
	GetLocationById(ctx context.Context, id int) (*Location, error)
	GetLocationsList(ctx context.Context, sortOptions middleware.SortOptions) ([]Location, error)
}

type serviceLocation struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo   RepositoryLocation
	logger logging.Logger
}

func NewLocationService(repo RepositoryLocation, logger logging.Logger) *serviceLocation {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceLocation{repo: repo, logger: logger}
}

func (sl *serviceLocation) CreateLocation(ctx context.Context, location *Location) (int, error) {
	newLocation, err := sl.repo.CreateLocation(ctx, location)

	if err != nil {
		sl.logger.Errorf("Failed to CreateLocation in service %s", err.Error())
		return 0, err
	}

	return newLocation, nil
}

func (sl *serviceLocation) GetLocationById(ctx context.Context, id int) (*Location, error) {

	locationById, err := sl.repo.GetLocationById(ctx, id)

	if err != nil {
		//sl.logger.Error("Failed to GetLocationById in service", zap.Error(err))
		return nil, err
	}

	return locationById, err

}

func (sl *serviceLocation) GetLocationsList(ctx context.Context, sortOptions middleware.SortOptions) ([]Location, error) {
	locationsList, err := sl.repo.GetLocationsList(ctx, sortOptions)

	if err != nil {
		sl.logger.Errorf("Failed to GetLocationsList in service %s", err.Error())
		return nil, err
	}

	return locationsList, err
}
