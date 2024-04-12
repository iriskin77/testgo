package users

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/iriskin77/testgo/constants"
	"github.com/iriskin77/testgo/pkg/logging"
)

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

	return fmt.Sprintf("%x", hash.Sum([]byte(constants.SaltPassword)))
}

// JWT Auth
type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func (su *serviceUser) GenerateToken(ctx context.Context, username, password string) (string, error) {
	userId, err := su.repo.GetUserByUsernamePassword(ctx, username, GeneratePasswordHash(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(constants.TimeTokenExpire).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userId,
	})

	return token.SignedString([]byte(constants.SigningKey))
}

func GetUserFromToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(constants.SigningKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	fmt.Println("ParseToken", token)
	return claims.UserId, nil
}
