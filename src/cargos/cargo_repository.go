package cargos

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	cargoTable = "cargos"
)

type RepositoryCargo interface {
	CreateCargo(ctx context.Context, cargo *models.CargoRequest) (int, error)
}

type CargoDB struct {
	db *pgxpool.Pool
}

func NewCargoDB(db *pgxpool.Pool) *CargoDB {
	return &CargoDB{db: db}
}

func (cr *CargoDB) CreateCargo(ctx context.Context, cargo *models.CargoRequest) (int, error) {

	var pickUpId int

	queryPickUp := `SELECT id FROM locations WHERE zip = $1`

	if err := cr.db.QueryRow(ctx, queryPickUp, cargo.Zip_pickup).Scan(
		&pickUpId); err != nil {
		return 0, err
	}

	var deliveryId int

	queryDelivery := `SELECT id FROM locations WHERE zip = $1`

	if err := cr.db.QueryRow(ctx, queryDelivery, cargo.Zip_delivery).Scan(
		&deliveryId); err != nil {
		return 0, err
	}

	var cargoId int

	query := fmt.Sprintf(`INSERT INTO %s (cargo_name, weight, description, pick_up_location_id, delivery_location_id) 
	                      VALUES ($1, $2, $3, $4, $5) 
						  RETURNING id`, cargoTable)

	if err := cr.db.QueryRow(ctx, query,
		cargo.Cargo_name,
		cargo.Weight,
		cargo.Description,
		pickUpId,
		deliveryId,
	).Scan(&cargoId); err != nil {
		return 0, err
	}

	return cargoId, nil
}

// CREATE TABLE cargos (
//     id serial unique not null,
//     cargo_name varchar(255) not null,
//     weight INT NOT NULL,
//     description varchar(1024) not null,
//     pick_up_location_id int REFERENCES locations (id) ON DELETE SET NULL,
//     delivery_location_id int REFERENCES locations (id) ON DELETE SET NULL
// );
