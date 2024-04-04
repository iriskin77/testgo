package locations

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	locationsTable = "locations"
)

type RepositoryLocation interface {
	CreateLocation(ctx context.Context, location *models.Location) (int, error)
	GetLocationById(ctx context.Context, id int) (*models.Location, error)
	GetLocationsList(ctx context.Context) ([]models.Location, error)
}

type LocationDB struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewLocationDB(db *pgxpool.Pool, logger *zap.Logger) *LocationDB {
	return &LocationDB{db: db, logger: logger}
}

func (lc *LocationDB) CreateLocation(ctx context.Context, location *models.Location) (int, error) {

	var id int

	query := fmt.Sprintf("INSERT INTO %s (city, state, zip, latitude, longitude) VALUES ($1, $2, $3, $4, $5) RETURNING id", locationsTable)

	if err := lc.db.QueryRow(ctx, query, location.City, location.State, location.Zip, location.Latitude, location.Longitude).Scan(&id); err != nil {
		lc.logger.Error("Failed to create location", zap.Error(err))
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
		lc.logger.Error("Failed to get a location by id", zap.Error(err))
		return &locationById, err
	}

	return &locationById, nil

}

func (lc *LocationDB) GetLocationsList(ctx context.Context) ([]models.Location, error) {

	locationsList := make([]models.Location, 0)

	query := fmt.Sprintf("SELECT id, city, state, zip, latitude, longitude, created_at FROM %s", locationsTable)

	rowsLocations, err := lc.db.Query(ctx, query)

	if err != nil {
		lc.logger.Error("Failed to retrieve list locations from db", zap.Error(err))
		return nil, err
	}

	for rowsLocations.Next() {
		var loc models.Location

		err = rowsLocations.Scan(
			&loc.Id,
			&loc.City,
			&loc.State,
			&loc.Zip,
			&loc.Latitude,
			&loc.Longitude,
			&loc.Created_at,
		)

		if err != nil {
			lc.logger.Error("Failed to retrieve list locations from db", zap.Error(err))
			return nil, err
		}

		locationsList = append(locationsList, loc)
	}

	// if err = rowsLocations.Err(); err != nil {
	// 	return nil, err
	// }

	return locationsList, nil
}
