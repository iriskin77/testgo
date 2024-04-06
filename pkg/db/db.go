package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
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
// func NewPostgresDB(cfg ConfigDB) (*sqlx.DB, error) {
// 	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
// 		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
// 	if err != nil {
// 		return nil, err
// 	}
// 	// С помощью функции Ping проверяется подключение
// 	err = db.Ping()
// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

// type Client interface {
// 	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
// 	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
// 	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
// 	Begin(ctx context.Context) (pgx.Tx, error)
// }

func NewPostgresDB(ctx context.Context, cfg ConfigDB) (dbpool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	dbpool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbpool, nil

}
