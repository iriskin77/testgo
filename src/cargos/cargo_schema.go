package cargos

import "github.com/iriskin77/testgo/models"

// for CreateCar handler
type CargoRequest struct {
	Cargo_name   string `json:"cargo_name"`
	Zip_pickup   int    `json:"zip_pickup"`
	Zip_delivery int    `json:"zip_delivery"`
	Weight       int    `json:"weight"`
	Description  string `json:"description"`
}

type CargoCarsResponse struct {
	Id           int             `json:"id"`
	Cargo_name   string          `json:"cargo_name"`
	Pickup_loc   models.Location `json:"pickup_loc"`
	Delivery_loc models.Location `json:"delivery_loc"`
	Weight       int             `json:"weight"`
	Description  string          `json:"description"`
	Cars         []CarResponse   `json:"cars"`
}

type CarResponse struct {
	Unique_number string          `json:"unique_number"`
	Car_name      string          `json:"car_name"`
	Load_capacity int             `json:"load_capacity"`
	Zip           int             `json:"zip"`
	Car_location  models.Location `json:"car_location"`
}

type CargoUpdateRequest struct {
	Id          int    `json:"id"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
}
