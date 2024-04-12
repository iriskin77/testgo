package cargos

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/iriskin77/testgo/internal/locations"
)

type Cargo struct {
	Id                int    `json:"id"`
	Cargo_name        string `json:"cargo_name"`
	Weight            int    `json:"weight"`
	Description       string `json:"description"`
	Pick_up_location  int    `json:"pick_up_location"`
	Delivery_location int    `json:"delivery_location"`
}

type CargoCreateRequest struct {
	Cargo_name   string `json:"cargo_name"`
	Weight       int    `json:"weight"`
	Description  string `json:"description"`
	Zip_pickup   int    `json:"zip_pickup"`
	Zip_delivery int    `json:"zip_delivery"`
}

func (c *CargoCreateRequest) CreateCargoValidate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Cargo_name, validation.Required, validation.Length(5, 50)),
		validation.Field(&c.Weight, validation.Required),
		validation.Field(&c.Description, validation.Required, validation.Length(5, 50)),
		validation.Field(&c.Zip_pickup, validation.Required),
		validation.Field(&c.Zip_delivery, validation.Required))
}

type CargoCarsResponse struct {
	Id           int                `json:"id"`
	Cargo_name   string             `json:"cargo_name"`
	Pickup_loc   locations.Location `json:"pickup_loc"`
	Delivery_loc locations.Location `json:"delivery_loc"`
	Weight       int                `json:"weight"`
	Description  string             `json:"description"`
	Cars         []CarResponse      `json:"cars"`
}

type CarResponse struct {
	Unique_number string             `json:"unique_number"`
	Car_name      string             `json:"car_name"`
	Load_capacity int                `json:"load_capacity"`
	Zip           int                `json:"zip"`
	Car_location  locations.Location `json:"car_location"`
}

type CargoUpdateRequest struct {
	Id          int    `json:"id"`
	Weight      int    `json:"weight"`
	Description string `json:"description"`
}

func (c *CargoUpdateRequest) UpdateCargoValidate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Id, validation.Required),
		validation.Field(&c.Weight, validation.Required),
		validation.Field(&c.Description, validation.Required),
	)
}
