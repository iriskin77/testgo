package repository

import (
	"github.com/iriskin77/testgo/models"
	"github.com/iriskin77/testgo/src/repository/filerepository"
	"github.com/jmoiron/sqlx"
)

// Интерфейсы называются в зависимости от участков доменной зоны, за которую они отвечают
type File interface {
	UploadFile(*models.File) int
	DownloadFile(id int) error
}

type Car interface {
}

type Location interface {
}

type Cargo interface {
}

type Repository struct {
	File
}

// Конструктор сервисов
// Поскольку репозиторий должен работать с БД, то
func NewRepository(db *sqlx.DB) *Repository {
	// В файле репозитория инициализируем наш репозиторий в конструкторе
	return &Repository{File: filerepository.NewFileDB(db)}
}
