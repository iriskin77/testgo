package cars

import (
	"time"

	"github.com/iriskin77/testgo/models"
)

type CarUpdateRequest struct {
	Id            int    `json:"id"`
	Unique_number string `json:"unique_number"`
	Car_name      string `json:"car_name"`
	Load_capacity int    `json:"load_capacity"`
	Zip           int    `json:"zip"`
}

type CarUpdateResponse struct {
	Id               int               `json:"id"`
	Unique_number    string            `json:"unique_number"`
	Car_name         string            `json:"car_name"`
	Load_capacity    int               `json:"load_capacity"`
	Created_at       time.Time         `json:"created_at"`
	New_car_location []models.Location `json:"car_location"`
}
