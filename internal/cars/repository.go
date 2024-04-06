package cars

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	carTable = "cars"
)

type RepositoryCar interface {
	CreateCar(ctx context.Context, car *CarRequest) (int, error)
	UpdateCarById(ctx context.Context, carUpdate *CarUpdateRequest) (int, error)
}

type CarDB struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

func NewCarDB(db *pgxpool.Pool, logger logging.Logger) *CarDB {
	return &CarDB{
		db:     db,
		logger: logger,
	}
}

func (c *CarDB) CreateCar(ctx context.Context, car *CarRequest) (int, error) {

	var locId int

	query_loc := "SELECT id FROM locations WHERE zip = $1"

	if err := c.db.QueryRow(ctx, query_loc, car.Zip).Scan(
		&locId); err != nil {
		return 0, err
	}

	var carId int

	query := fmt.Sprintf("INSERT INTO %s (unique_number, car_name, load_capacity, car_location_id) VALUES ($1, $2, $3, $4) RETURNING id", carTable)

	if err := c.db.QueryRow(ctx, query,
		car.Unique_number,
		car.Car_name,
		car.Load_capacity,
		locId).Scan(&carId); err != nil {
		return 0, err
	}

	return carId, nil
}

func (c *CarDB) UpdateCarById(ctx context.Context, carUpdate *CarUpdateRequest) (int, error) {

	var newLocationId int
	var carUpdatedId int

	fmt.Println(carUpdate.Zip)

	queryZip := `SELECT id
	             FROM locations
				 WHERE zip = $1`

	if err := c.db.QueryRow(ctx, queryZip,
		carUpdate.Zip).Scan(&newLocationId); err != nil {
		c.logger.Errorf("Failed to get location by zip from DB %s", err.Error())
		return 0, err
	}

	fmt.Println(newLocationId)

	queryCarUpdate := `UPDATE cars
	                   SET unique_number = $1, car_name = $2, load_capacity = $3, car_location_id = $4 
					   WHERE id = $5
					   RETURNING id`

	if err := c.db.QueryRow(ctx, queryCarUpdate,
		&carUpdate.Unique_number,
		&carUpdate.Car_name,
		&carUpdate.Load_capacity,
		newLocationId,
		&carUpdate.Id).Scan(&carUpdatedId); err != nil {
		return 0, err
	}

	return carUpdatedId, nil

}
