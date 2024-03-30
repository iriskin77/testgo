package models

type Car struct {
	Id            int    `json:"id"`
	Unique_number string `json:"unique_number"`
	Car_name      string `json:"car_name"`
	Load_capacity int    `json:"load_capacity"`
	Created_at    string `json:"created_at"`
	Car_location  string `json:"car_location"`
}
