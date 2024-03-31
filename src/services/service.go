package services

import (
	"github.com/iriskin77/testgo/models"
	"github.com/iriskin77/testgo/src/repository"
)

type File interface {
	UploadFile(*models.File) int
	DownloadFile(id int) (*models.File, error)
}

type Car interface {
}

type Location interface {
}

type Cargo interface {
}

type Service struct {
	File
}

// Конструктор сервисов. Сервисы будут передавать данные из хэндлера ниже, на уровень репозитория, поэтому нужен указатель
// на структуру репозитория (репозиторий коннектиться к БД)

func NewService(repository *repository.Repository) *Service {
	return &Service{
		File: NewFileService(repository.File),
	}
}
