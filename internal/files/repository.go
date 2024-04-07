package files

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	filesTable    = "file"
	locationTable = "locations"
)

type RepositoryFile interface {
	UploadFile(ctx context.Context, file *File) (int, error)
	DownloadFile(ctx context.Context, id int) (*File, error)
	BulkInsertLocations(ctx context.Context, id int) (int, error)
}

type FileDB struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

func NewFileDB(db *pgxpool.Pool, logger logging.Logger) *FileDB {
	return &FileDB{
		db:     db,
		logger: logger,
	}
}

func (f *FileDB) UploadFile(ctx context.Context, file *File) (int, error) {
	var id int

	query := fmt.Sprintf("INSERT INTO %s (name, file_path) VALUES ($1, $2) RETURNING id", filesTable)
	if err := f.db.QueryRow(ctx, query, file.Name, file.File_path).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (f *FileDB) DownloadFile(ctx context.Context, id int) (*File, error) {

	var fileById File

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", filesTable)

	err := f.db.QueryRow(ctx, query, id).Scan(&fileById)
	if err != nil {
		return &fileById, err
	}

	return &fileById, nil

}

func (f *FileDB) BulkInsertLocations(ctx context.Context, id int) (int, error) {

	var filePath string

	queryGetFile := fmt.Sprintf("SELECT file_path FROM %s WHERE id = $1", filesTable)

	if err := f.db.QueryRow(ctx, queryGetFile, id).Scan(
		&filePath,
	); err != nil {
		f.logger.Errorf("Failed to get a location by id %s", err.Error())
		return 0, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','

	rowsInsert := make([][]interface{}, 0)

	for {

		listCsvRows := make([]interface{}, 0)
		record, e := reader.Read()
		if e != nil {
			fmt.Println(1)
			break
		}
		state := record[4]
		city := record[3]
		zip, _ := strconv.Atoi(record[0])
		latitude, _ := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		longitude, _ := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)

		listCsvRows = append(listCsvRows, state, city, zip, latitude, longitude)
		rowsInsert = append(rowsInsert, listCsvRows)

	}

	rows := [][]interface{}{
		{"Hadley", "Massachusetts", int32(601), float32(44.44), float32(44.44)},
	}

	fmt.Println(rows)

	CopyFile, err := f.db.CopyFrom(
		ctx,
		pgx.Identifier{"locations"},
		[]string{"city", "state", "zip", "latitude", "longitude"},
		pgx.CopyFromRows(rowsInsert),
	)

	return int(CopyFile), nil

}
