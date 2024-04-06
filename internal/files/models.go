package files

import "time"

type File struct {
	Id         int       `json:"id"`
	Name       string    `json:"name"`
	File_path  string    `json:"file_path"`
	Created_at time.Time `json:"created_at"`
}
