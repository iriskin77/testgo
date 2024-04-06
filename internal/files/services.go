package files

import (
	"context"

	"go.uber.org/zap"
)

type ServiceFile interface {
	UploadFile(ctx context.Context, file *File) (int, error)
	DownloadFile(ctx context.Context, id int) (*File, error)
}

type serviceFile struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo   RepositoryFile
	logger *zap.Logger
}

func NewFileService(repo RepositoryFile, logger *zap.Logger) *serviceFile {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceFile{
		repo:   repo,
		logger: logger,
	}
}

func (s *serviceFile) UploadFile(ctx context.Context, file *File) (int, error) {
	fileId, _ := s.repo.UploadFile(ctx, file)
	return fileId, nil
}

func (s *serviceFile) DownloadFile(ctx context.Context, id int) (*File, error) {
	file, err := s.repo.DownloadFile(ctx, id)

	if err != nil {
		return nil, err
	}

	return file, nil
}