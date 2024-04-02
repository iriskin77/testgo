package models

import "time"

type File struct {
	Id         int
	Name       string
	File_path  string
	Created_at time.Time
}
