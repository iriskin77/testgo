package locations

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/iriskin77/testgo/internal/middleware"
	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	locationsTable = "locations"
)

type RepositoryLocation interface {
	CreateLocation(ctx context.Context, location *Location) (int, error)
	GetLocationById(ctx context.Context, id int) (*Location, error)
	GetLocationsList(ctx context.Context, sortOptions middleware.SortOptions) ([]Location, error)
}

type LocationDB struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

func NewLocationDB(db *pgxpool.Pool, logger logging.Logger) *LocationDB {
	return &LocationDB{db: db, logger: logger}
}

func (lc *LocationDB) CreateLocation(ctx context.Context, location *Location) (int, error) {

	var id int

	query := fmt.Sprintf("INSERT INTO %s (city, state, zip, latitude, longitude) VALUES ($1, $2, $3, $4, $5) RETURNING id", locationsTable)

	if err := lc.db.QueryRow(ctx, query, location.City, location.State, location.Zip, location.Latitude, location.Longitude).Scan(&id); err != nil {
		lc.logger.Errorf("Failed to create location %s", err.Error())
		return 0, err
	}

	return id, nil

}

func (lc *LocationDB) GetLocationById(ctx context.Context, id int) (*Location, error) {

	var locationById Location

	query := fmt.Sprintf("SELECT id, city, state, zip, latitude, longitude, created_at FROM %s WHERE id = $1", locationsTable)

	if err := lc.db.QueryRow(ctx, query, id).Scan(
		&locationById.Id,
		&locationById.City,
		&locationById.State,
		&locationById.Zip,
		&locationById.Latitude,
		&locationById.Longitude,
		&locationById.Created_at); err != nil {
		lc.logger.Errorf("Failed to get a location by id %s", err.Error())
		return &locationById, err
	}

	return &locationById, nil

}

func (lc *LocationDB) GetLocationsList(ctx context.Context, sortOptions middleware.SortOptions) ([]Location, error) {

	locationsList := make([]Location, 0)

	qb := sq.Select("id", "city", "state", "zip", "latitude", "longitude", "created_at").
		From(locationsTable)

	if sortOptions.Field != "" && sortOptions.Order != "" {
		qb.OrderBy(fmt.Sprintf("%s %s", sortOptions.Field, sortOptions.Order))
	}

	query, sortOpt, err := qb.ToSql()

	if err != nil {
		return locationsList, err
	}

	//query := fmt.Sprintf("SELECT id, city, state, zip, latitude, longitude, created_at FROM %s", locationsTable)

	rowsLocations, err := lc.db.Query(ctx, query, sortOpt...)

	if err != nil {
		lc.logger.Errorf("Failed to retrieve list locations from db %s", err.Error())
		return nil, err
	}

	for rowsLocations.Next() {
		var loc Location

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
			lc.logger.Errorf("Failed to retrieve list locations from db %s", err.Error())
			return nil, err
		}

		locationsList = append(locationsList, loc)
	}

	return locationsList, nil
}
