package models

type Cargo struct {
	Id                int    `json:"id"`
	Cargo_name        string `json:"cargo_name"`
	Weight            int    `json:"weight"`
	Description       string `json:"description"`
	Pick_up_location  int    `json:"pick_up_location"`
	Delivery_location int    `json:"delivery_location"`
}

// for CreateCar handler
type CargoRequest struct {
	Cargo_name   string `json:"cargo_name"`
	Zip_pickup   int    `json:"zip_pickup"`
	Zip_delivery int    `json:"zip_delivery"`
	Weight       int    `json:"weight"`
	Description  string `json:"description"`
}
