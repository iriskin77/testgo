package users

import (
	"context"
	"fmt"

	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable = "users"
)

type RepositoryUser interface {
	CreateUser(ctx context.Context, newUser *User) (int, error)
	GetUserByUsernamePassword(ctx context.Context, username, password string) (int, error)
}

type UserDB struct {
	db     *pgxpool.Pool
	logger logging.Logger
}

func NewUserDB(db *pgxpool.Pool, logger logging.Logger) *UserDB {
	return &UserDB{db: db, logger: logger}
}

func (ru *UserDB) CreateUser(ctx context.Context, newUser *User) (int, error) {

	var newUserId int

	query := fmt.Sprintf("INSERT INTO %s (name, surname, age, email, username, password_hash) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", usersTable)

	if err := ru.db.QueryRow(ctx, query,
		newUser.Name,
		newUser.Surname,
		newUser.Age,
		newUser.Email,
		newUser.Username,
		newUser.Password_hash,
	).Scan(&newUserId); err != nil {
		return 0, err
	}

	return newUserId, nil
}

func (ru *UserDB) GetUserByUsernamePassword(ctx context.Context, username, password string) (int, error) {

	var userId int

	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)

	if err := ru.db.QueryRow(ctx, query,
		username,
		password,
	).Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}
