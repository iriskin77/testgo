package cargos

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/internal/locations"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

const (
	cargoTable = "cargos"
)

type RepositoryCargo interface {
	CreateCargo(ctx context.Context, cargo *CargoRequest) (int, error)
	GetCargoCars(ctx context.Context, id int) (*CargoCarsResponse, error)
	GetListCargos(ctx context.Context) ([]interface{}, error)
	UpdateCargoById(ctx context.Context, cargoUpdate *CargoUpdateRequest) (int, error)
}

type CargoDB struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewCargoDB(db *pgxpool.Pool, logger *zap.Logger) *CargoDB {
	return &CargoDB{
		db:     db,
		logger: logger,
	}
}

func (cr *CargoDB) CreateCargo(ctx context.Context, cargo *CargoRequest) (int, error) {

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

func (cr *CargoDB) GetCargoCars(ctx context.Context, id int) (*CargoCarsResponse, error) {

	// final response with the choosen cargo and cars
	var cargoCars CargoCarsResponse
	// cars list that are closest to the choosen cargo
	var cars []CarResponse

	// locations for cargo
	var CargoPickUpLocation locations.Location
	var CargoDeliveryLocation locations.Location

	var cargoPickUpId int
	var cargoDeliveryId int

	queryCargo := `SELECT cargo_name, weight, description, pick_up_location_id, delivery_location_id
	               FROM cargos 
				   WHERE id = $1`

	if err := cr.db.QueryRow(ctx, queryCargo, id).Scan(
		&cargoCars.Cargo_name,
		&cargoCars.Weight,
		&cargoCars.Description,
		&cargoPickUpId,
		&cargoDeliveryId,
	); err != nil {
		return &cargoCars, err
	}

	fmt.Println(cargoCars)

	queryPickUpLocation := `SELECT id, city, state, zip, latitude, longitude
	                        FROM locations
							WHERE id = $1`

	if err := cr.db.QueryRow(ctx, queryPickUpLocation, cargoPickUpId).Scan(
		&CargoPickUpLocation.Id,
		&CargoPickUpLocation.City,
		&CargoPickUpLocation.State,
		&CargoPickUpLocation.Zip,
		&CargoPickUpLocation.Latitude,
		&CargoPickUpLocation.Longitude); err != nil {
		return &cargoCars, err
	}

	fmt.Println(CargoPickUpLocation)

	queryDeliveryLocation := `SELECT id, city, state, zip, latitude, longitude
	                          FROM locations
							  WHERE id = $1`

	if err := cr.db.QueryRow(ctx, queryDeliveryLocation, cargoDeliveryId).Scan(
		&CargoDeliveryLocation.Id,
		&CargoDeliveryLocation.City,
		&CargoDeliveryLocation.State,
		&CargoDeliveryLocation.Zip,
		&CargoDeliveryLocation.Latitude,
		&CargoDeliveryLocation.Longitude); err != nil {
		return &cargoCars, err
	}

	fmt.Println(CargoDeliveryLocation)

	//Получаем машины и локации

	queryCars := `SELECT cars.unique_number, cars.car_name, cars.load_capacity,
	                     locations.id, locations.city, locations.state, locations.zip, locations.latitude, locations.longitude
	                          
				  FROM cars INNER JOIN locations ON cars.id = locations.id`

	rowsCars, err := cr.db.Query(ctx, queryCars)

	if err != nil {
		return nil, err
	}

	fmt.Print(rowsCars.Next())

	for rowsCars.Next() {
		var car CarResponse
		//var loc models.Location

		err = rowsCars.Scan(
			&car.Unique_number,
			&car.Car_name,
			&car.Load_capacity,
			&car.Car_location.Id,
			&car.Car_location.City,
			&car.Car_location.State,
			&car.Car_location.Zip,
			&car.Car_location.Latitude,
			&car.Car_location.Longitude,
		)

		if err != nil {
			return nil, err
		}

		cars = append(cars, car)
	}

	if err := rowsCars.Err(); err != nil {
		return nil, err
	}

	fmt.Println(cars)

	cargoCars.Pickup_loc = CargoPickUpLocation
	cargoCars.Delivery_loc = CargoDeliveryLocation
	cargoCars.Cars = cars

	return &cargoCars, nil
}

func (cr *CargoDB) GetListCargos(ctx context.Context) ([]interface{}, error) {

	var CargoCarsResp []interface{}

	queryCargos := `SELECT id 
				    FROM cargos`

	rowsCargos, err := cr.db.Query(ctx, queryCargos)

	if err != nil {
		return nil, err
	}

	for rowsCargos.Next() {
		var cargoId int

		err = rowsCargos.Scan(
			&cargoId,
		)

		if err != nil {
			return nil, err
		}

		cargoItem, err := cr.GetCargoCars(ctx, cargoId)

		if err != nil {
			return nil, err
		}

		CargoCarsResp = append(CargoCarsResp, cargoItem)

	}

	return CargoCarsResp, nil

}

func (cr *CargoDB) UpdateCargoById(ctx context.Context, cargoUpdate *CargoUpdateRequest) (int, error) {

	var cargoUpdatedId int

	queryUpdateCargo := `UPDATE cargos
	                     SET weight = $1, description = $2
						 WHERE id = $1
						 RETURNING id`

	if err := cr.db.QueryRow(ctx, queryUpdateCargo,
		&cargoUpdate.Weight,
		&cargoUpdate.Description,
		&cargoUpdate.Id).Scan(&cargoUpdatedId); err != nil {
		cr.logger.Error("Failed to update cargo in DB", zap.Error(err))
		return 0, err
	}

	return cargoUpdatedId, nil
}
