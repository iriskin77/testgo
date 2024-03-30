package file_service

import (
	"github.com/iriskin77/testgo/models"
	"github.com/iriskin77/testgo/src/repository"
)

type ServiceFile struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo repository.File
}

func NewFileService(repo repository.File) *ServiceFile {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &ServiceFile{repo: repo}
}

func (s *ServiceFile) UploadFile(*models.File) int {
	return 1
}

func (s *ServiceFile) DownloadFile(id int) error {
	return nil
}
