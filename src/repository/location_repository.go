package repository

import (
	"fmt"

	"github.com/iriskin77/testgo/models"
	"github.com/jmoiron/sqlx"
)

const (
	locationsTable = "location"
)

type LocationDB struct {
	db *sqlx.DB
}

func NewLocationDB(db *sqlx.DB) *LocationDB {
	return &LocationDB{db: db}
}

func (l *LocationDB) InsertFileToDB(fileId int) {

	var fileById models.File

	query := fmt.Sprintf("SELECT * FROM file WHERE id = 2")
	err := l.db.Get(&fileById, query, fileId)

	queryq := fmt.Sprintf("COPY persons(zip, latitude, longitude, city, state) FROM '/home/abc/Рабочий стол/uszips.csv' DELIMITER ',' CSV HEADER;")
	_, err = l.db.Exec(queryq)

	if err != nil {
		fmt.Println(err)
	}

	return

}

// city VARCHAR(255) NOT NULL,
//     state VARCHAR(255) NOT NULL,
//     zip INT NOT NULL,
//     latitude FLOAT NOT NULL,
//     longitude FLOAT NOT NULL,

// COPY persons(first_name, last_name, dob, email)
// FROM 'C:\sampledb\persons.csv'
// DELIMITER ','
// CSV HEADER;
