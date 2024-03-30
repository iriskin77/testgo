package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type ConfigDB struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

// Функция NewPostgresDB возвращает указатель на структуру sqlxDB
func NewPostgresDB(cfg ConfigDB) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	// С помощью функции Ping проверяется подключение
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
