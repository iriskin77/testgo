package users

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/iriskin77/testgo/pkg/logging"
)

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type ServiceUser interface {
	CreateUser(ctx context.Context, newUser *User) (int, error)
	GenerateToken(ctx context.Context, username, password string) (string, error)
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

func (su *serviceUser) GenerateToken(ctx context.Context, username, password string) (string, error) {
	userId, err := su.repo.GetUserByUsernamePassword(ctx, username, GeneratePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(12 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	return token.SignedString([]byte(signingKey))
}
