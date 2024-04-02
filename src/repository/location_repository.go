package repository

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	locationsTable = "locations"
)

type LocationDB struct {
	db *pgxpool.Pool
}

func NewLocationDB(db *pgxpool.Pool) *LocationDB {
	return &LocationDB{db: db}
}

func (lc *LocationDB) CreateLocation(ctx context.Context, location *models.Location) (int, error) {

	var id int

	query := fmt.Sprintf("INSERT INTO %s (city, state, zip, latitude, longitude) VALUES ($1, $2, $3, $4, $5) RETURNING id", locationsTable)

	if err := lc.db.QueryRow(ctx, query, location.City, location.State, location.Zip, location.Latitude, location.Longitude).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil

}

func (lc *LocationDB) GetLocationById(ctx context.Context, id int) (*models.Location, error) {

	var locationById models.Location

	query := fmt.Sprintf("SELECT id, city, state, zip, latitude, longitude, created_at FROM %s WHERE id = $1", locationsTable)

	if err := lc.db.QueryRow(ctx, query, id).Scan(
		&locationById.Id,
		&locationById.City,
		&locationById.State,
		&locationById.Zip,
		&locationById.Latitude,
		&locationById.Longitude,
		&locationById.Created_at); err != nil {
		return &locationById, err
	}

	return &locationById, nil

}

func (lc *LocationDB) GetLocationsList(ctx context.Context) ([]models.Location, error) {
	return nil, nil
}
