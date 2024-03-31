package repository

import (
	"github.com/iriskin77/testgo/models"
	"github.com/jmoiron/sqlx"
)

// Интерфейсы называются в зависимости от участков доменной зоны, за которую они отвечают
type File interface {
	UploadFile(*models.File) int
	DownloadFile(id int) (*models.File, error)
}

type Car interface {
}

type Location interface {
	InsertFileToDB(fileId int)
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
func NewRepository(db *sqlx.DB) *Repository {
	// В файле репозитория инициализируем наш репозиторий в конструкторе
	return &Repository{File: NewFileDB(db),
		Location: NewLocationDB(db)}
}
