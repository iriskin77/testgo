package files

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	filesTable    = "file"
	locationTable = "locations"
)

type RepositoryFile interface {
	UploadFile(ctx context.Context, file *models.File) (int, error)
	DownloadFile(ctx context.Context, id int) (*models.File, error)
}

type FileDB struct {
	db *pgxpool.Pool
}

func NewFileDB(db *pgxpool.Pool) *FileDB {
	return &FileDB{db: db}
}

func (f *FileDB) UploadFile(ctx context.Context, file *models.File) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, file_path) VALUES ($1, $2) RETURNING id", filesTable)
	if err := f.db.QueryRow(ctx, query, file.Name, file.File_path).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (f *FileDB) DownloadFile(ctx context.Context, id int) (*models.File, error) {

	var fileById models.File

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", filesTable)

	err := f.db.QueryRow(ctx, query, id).Scan(&fileById)
	if err != nil {
		return &fileById, err
	}

	return &fileById, nil

}
