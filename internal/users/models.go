package users

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	Id            int       `json:"id"`
	Name          string    `json:"name"`
	Surname       string    `json:"surname"`
	Age           int       `json:"age"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password_hash string    `json:"password"`
	Updated_at    time.Time `json:"updated at"`
	Created_at    time.Time `json:"created_at"`
}

func (u *User) CreateUserValidate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Name, validation.Required, validation.Length(5, 50)),
		validation.Field(&u.Surname, validation.Required, validation.Length(5, 50)),
		validation.Field(&u.Age, validation.Required),
		validation.Field(&u.Username, validation.Required),
		validation.Field(&u.Email, validation.Required),
		validation.Field(&u.Password_hash, validation.Required))
}

// validation.Match(regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`))
// validation.Match(regexp.MustCompile("^[0-9A-Za-z]+$")
