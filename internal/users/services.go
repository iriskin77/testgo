package users

import (
	"context"
	"crypto/sha1"
	"fmt"

	"github.com/iriskin77/testgo/pkg/logging"
)

const (
	salt = "hjqrhjqw124617ajfhajs"
)

type ServiceUser interface {
	CreateUser(ctx context.Context, newUser *User) (int, error)
}

type serviceUser struct {
	// создаем структуру, которая принимает репозиторий для работы с БД
	repo   RepositoryUser
	logger logging.Logger
}

func NewUserService(repo RepositoryUser, logger logging.Logger) *serviceUser {
	// Конструктор: принимает репозиторий, возваращает сервис с репозиторием
	return &serviceUser{repo: repo, logger: logger}
}

func (su *serviceUser) CreateUser(ctx context.Context, newUser *User) (int, error) {
	newUser.Password_hash = GeneratePasswordHash(newUser.Password_hash)
	newUserId, err := su.repo.CreateUser(ctx, newUser)

	if err != nil {
		return 0, err
	}

	return newUserId, nil

}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
