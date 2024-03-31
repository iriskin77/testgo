package services

import "github.com/iriskin77/testgo/src/repository"

type ServiceLocation struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo repository.Location
}

func NewLocationService(repo repository.Location) *ServiceLocation {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &ServiceLocation{repo: repo}
}

func (sl *ServiceLocation) InsertFileToDB(fileId int) {
	sl.repo.InsertFileToDB(fileId)

}
