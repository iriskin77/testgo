package cargos

import (
	"context"

	"github.com/iriskin77/testgo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RepositoryCargo interface {
	CreateCar(ctx context.Context, car *models.Cargo) (int, error)
}

type CargoDB struct {
	db *pgxpool.Pool
}

func NewCargoDB(db *pgxpool.Pool) *CargoDB {
	return &CargoDB{db: db}
}

func (cr *CargoDB) CreateCar(ctx context.Context, car *models.Cargo) (int, error) {
	return 1, nil
}
