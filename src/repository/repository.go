package repository

import (
	"context"

	"github.com/iriskin77/testgo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Интерфейсы называются в зависимости от участков доменной зоны, за которую они отвечают
type File interface {
	UploadFile(ctx context.Context, file *models.File) (int, error)
	DownloadFile(ctx context.Context, id int) (*models.File, error)
}

type Car interface {
	CreateCar(ctx context.Context, car *models.CarRequest) (int, error)
}

type Location interface {
	CreateLocation(ctx context.Context, location *models.Location) (int, error)
	GetLocationById(ctx context.Context, id int) (*models.Location, error)
	GetLocationsList(ctx context.Context) ([]models.Location, error)
}

type Cargo interface {
}

type Repository struct {
	File
	Location
	Car
	Cargo
}

// Конструктор сервисов
// Поскольку репозиторий должен работать с БД, то
func NewRepository(db *pgxpool.Pool) *Repository {
	// В файле репозитория инициализируем наш репозиторий в конструкторе
	return &Repository{File: NewFileDB(db),
		Location: NewLocationDB(db),
		Car:      NewCarDB(db),
	}
}
