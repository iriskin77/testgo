package repository

import (
	"fmt"

	"github.com/iriskin77/testgo/models"
	"github.com/jmoiron/sqlx"
)

const (
	filesTable = "file"
)

type FileDB struct {
	db *sqlx.DB
}

func NewFileDB(db *sqlx.DB) *FileDB {
	return &FileDB{db: db}
}

func (f *FileDB) UploadFile(file *models.File) int {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, file) VALUES ($1, $2) RETURNING id", filesTable)
	row := f.db.QueryRow(query, file.Name, file.File)
	if err := row.Scan(&id); err != nil {
		return 0
	}

	return id
}

func (f *FileDB) DownloadFile(id int) error {
	return nil
}
