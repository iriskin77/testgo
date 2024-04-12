package locations

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Location struct {
	Id         int       `json:"id"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Zip        int       `json:"zip"`
	Latitude   float32   `json:"latitude"`
	Longitude  float32   `json:"longitude"`
	Created_at time.Time `json:"created_at"`
}

func (l *Location) CreateLocationValidate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.City, validation.Required, validation.Length(5, 50)),
		validation.Field(&l.State, validation.Required, validation.Length(5, 50)),
		validation.Field(&l.Zip, validation.Required),
		validation.Field(&l.Latitude, validation.Required),
		validation.Field(&l.Longitude, validation.Required))
}
