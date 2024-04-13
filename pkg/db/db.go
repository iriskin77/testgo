package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/iriskin77/testgo/configs"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresDB(ctx context.Context, cfg configs.ConfigPostgres) (dbpool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.NameDB)
	dbpool, err = pgxpool.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	return dbpool, nil

}
