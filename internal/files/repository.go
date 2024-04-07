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
		fmt.Println(1)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ','
	//batch := &pgx.Batch{}

	//queryInsert := "INSERT INTO locations (city, state, zip, latitude, longitude) VALUES (@city, @state, @zip, @latitude, @longitude)"
	//csv_rows := make([][]interface{}, 0)
	//locationsList := make([]Location, 0)
	rows1 := make([][]interface{}, 0)
	for {
		lst := make([]interface{}, 0)
		record, e := reader.Read()
		if e != nil {
			fmt.Println(1)
			break
		}
		//fmt.Println(record)

		// fmt.Println(city)
		state := record[4]
		city := record[3]
		// fmt.Println(state)
		zip, _ := strconv.Atoi(record[0])
		// fmt.Println(zip)
		latitude, _ := strconv.ParseFloat(strings.TrimSpace(record[1]), 64)
		// fmt.Println(latitude)
		longitude, _ := strconv.ParseFloat(strings.TrimSpace(record[2]), 64)
		//feetFloat, _ := strconv.ParseFloat(strings.TrimSpace(variable), 64)
		//zip_int, _ := strconv.Atoi(zip)
		// lat_float, _ := strconv.Atoi(latitude)
		// long_float, _ := strconv.Atoi(longitude)

		lst = append(lst, state, city, zip, latitude, longitude)
		rows1 = append(rows1, lst)

		// fmt.Println(city)
		// fmt.Println(state)
		// fmt.Println(zip)
		// fmt.Println(latitude)
		// fmt.Println(longitude)

		//queryInsert := fmt.Sprintf(`INSERT INTO locations (city, state, zip, latitude, longitude)
		// VALUES (%s, %s, %s, %s, %s)`, city, state, zip, latitude, longitude)

		// args := pgx.NamedArgs{
		// 	"@city":     city,
		// 	"@state":    state,
		// 	"@zip":      zip,
		// 	"latitude":  12,
		// 	"longitude": 12,
		// }

		//fmt.Println(args)
		//batch.Queue(queryInsert)
	}

	// for k, v := range rows1 {
	// 	fmt.Println(k, v)
	// }

	rows := [][]interface{}{
		{"Hadley", "Massachusetts", int32(601), float32(44.44), float32(44.44)},
	}

	//c := ["Hadley", "Massachusetts", int32(601), float32(44.44), float32(44.44)]
	//rows = append(rows, ["Hadley", "Massachusetts", int32(601), float32(44.44), float32(44.44)])

	fmt.Println(rows)

	copyCount, err := f.db.CopyFrom(
		ctx,
		pgx.Identifier{"locations"},
		[]string{"city", "state", "zip", "latitude", "longitude"},
		pgx.CopyFromRows(rows1),
	)

	fmt.Println(copyCount)

	// results := f.db.SendBatch(ctx, batch)
	// defer results.Close()
	// // results.Exec()
	// // //defer results.Close()
	// // f.db.CopyFrom()
	// fmt.Println(results)

	// for i := 0; 5 > i; i++ {
	// 	q, err := results.Exec()
	// 	fmt.Println(q)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	// results.Close()
	return 1, nil

}
