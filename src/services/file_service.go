package services

import (
	"context"

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

func (s *ServiceFile) UploadFile(ctx context.Context, file *models.File) (int, error) {
	fileId, _ := s.repo.UploadFile(ctx, file)
	return fileId, nil
}

func (s *ServiceFile) DownloadFile(ctx context.Context, id int) (*models.File, error) {
	file, err := s.repo.DownloadFile(ctx, id)

	if err != nil {
		return nil, err
	}

	return file, nil
}
