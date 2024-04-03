package cars

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	carTable = "cars"
)

type RepositoryCar interface {
	CreateCar(ctx context.Context, car *models.CarRequest) (int, error)
}

type CarDB struct {
	db *pgxpool.Pool
}

func NewCarDB(db *pgxpool.Pool) *CarDB {
	return &CarDB{db: db}
}

func (c *CarDB) CreateCar(ctx context.Context, car *models.CarRequest) (int, error) {

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
