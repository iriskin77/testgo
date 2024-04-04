package cargos

import (
	"context"
)

type ServiceCar interface {
	CreateCargo(ctx context.Context, cargo *CargoRequest) (int, error)
	GetCargoCars(ctx context.Context, id int) (*CargoCarsResponse, error)
	GetListCargos(ctx context.Context) ([]interface{}, error)
}

type serviceCargo struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo RepositoryCargo
}

func NewCargoService(repo RepositoryCargo) *serviceCargo {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceCargo{repo: repo}
}

func (cr *serviceCargo) CreateCargo(ctx context.Context, cargo *CargoRequest) (int, error) {
	newCarId, err := cr.repo.CreateCargo(ctx, cargo)

	if err != nil {
		return 0, err
	}

	return newCarId, nil
}

func (cr *serviceCargo) GetCargoCars(ctx context.Context, id int) (*CargoCarsResponse, error) {

	car, err := cr.repo.GetCargoCars(ctx, id)

	if err != nil {
		return nil, err
	}

	return car, nil
}

func (cr *serviceCargo) GetListCargos(ctx context.Context) ([]interface{}, error) {

	car, err := cr.repo.GetListCargos(ctx)

	if err != nil {
		return nil, err
	}

	return car, nil
}
