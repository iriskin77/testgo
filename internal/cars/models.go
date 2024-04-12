package cars

import (
	"time"

	"github.com/iriskin77/testgo/internal/locations"
)

// main model
type Car struct {
	Id              int       `json:"id"`
	Unique_number   string    `json:"unique_number"`
	Car_name        string    `json:"car_name"`
	Load_capacity   int       `json:"load_capacity"`
	Created_at      time.Time `json:"created_at"`
	Car_location_id int       `json:"car_location"`
}

// for CreateCar handler
type CarCreateRequest struct {
	Unique_number string `json:"unique_number"`
	Car_name      string `json:"car_name"`
	Load_capacity int    `json:"load_capacity"`
	Zip           int    `json:"zip"`
	//Car_location  string `json:"car_location"`
}

type CarUpdateRequest struct {
	Id            int    `json:"id"`
	Unique_number string `json:"unique_number"`
	Car_name      string `json:"car_name"`
	Load_capacity int    `json:"load_capacity"`
	Zip           int    `json:"zip"`
}

type CarUpdateResponse struct {
	Id               int                  `json:"id"`
	Unique_number    string               `json:"unique_number"`
	Car_name         string               `json:"car_name"`
	Load_capacity    int                  `json:"load_capacity"`
	Created_at       time.Time            `json:"created_at"`
	New_car_location []locations.Location `json:"car_location"`
}
