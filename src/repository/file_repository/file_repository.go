package file_repository

import (
	"github.com/iriskin77/testgo/models"
	"github.com/jmoiron/sqlx"
)

type FileDB struct {
	db *sqlx.DB
}

func NewFileDB(db *sqlx.DB) *FileDB {
	return &FileDB{db: db}
}

func (f *FileDB) UploadFile(*models.File) int {
	return 1
}

func (f *FileDB) DownloadFile(id int) error {
	return nil
}
