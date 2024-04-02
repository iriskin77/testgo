package services

import (
	"context"

	"github.com/iriskin77/testgo/models"
	"github.com/iriskin77/testgo/src/repository"
)

type File interface {
	UploadFile(ctx context.Context, file *models.File) (int, error)
	DownloadFile(ctx context.Context, id int) (*models.File, error)
}

type Car interface {
}

type Location interface {
	CreateLocation(ctx context.Context, location *models.Location) (int, error)
	GetLocationById(ctx context.Context, id int) (*models.Location, error)
	GetLocationsList(ctx context.Context) ([]models.Location, error)
}

type Cargo interface {
}

type Service struct {
	File
	Location
	Car
	Cargo
}

// Конструктор сервисов. Сервисы будут передавать данные из хэндлера ниже, на уровень репозитория, поэтому нужен указатель
// на структуру репозитория (репозиторий коннектиться к БД)

func NewService(repository *repository.Repository) *Service {
	return &Service{
		File:     NewFileService(repository.File),
		Location: NewLocationService(repository.Location),
	}
}
