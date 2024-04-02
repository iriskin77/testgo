package models

import "time"

type Location struct {
	Id         int       `json:"id"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	Zip        int       `json:"zip"`
	Latitude   float32   `json:"latitude"`
	Longitude  float32   `json:"longitude"`
	Created_at time.Time `json:"created_at"`
}

// type User struct {
// 	Id            int    `json:"id"`
// 	Name          string `json:"name"`
// 	Surname       string `json:"surname"`
// 	Age           int    `json:"age"`
// 	Password_hash string `json:"password_hash"`
// 	Email         string `json:"email"`
// 	Is_admin      bool   `json:"is_admin"`
// }
