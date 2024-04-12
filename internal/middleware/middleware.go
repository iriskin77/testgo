package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

// Middleware for sorting

const (
	AscSort           = "ASC"
	DescSort          = "DESC"
	OptionsContextKey = "sort_options"
	DefaulSortField   = "created_at"
	DefaultSortOrder  = "ASC"
)

type SortOptions struct {
	Field string
	Order string
}

func SortMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		sortBy := request.URL.Query().Get("sort_by")
		sortOrder := request.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = DefaulSortField
		}
		if sortOrder == "" {
			sortOrder = DefaultSortOrder
		} else {
			upperSortOrder := strings.ToUpper(sortOrder)
			if upperSortOrder != AscSort && upperSortOrder != DescSort {
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte("Нельзя так делать"))

			}
		}

		options := SortOptions{
			Field: sortBy,
			Order: sortOrder,
		}

		ctx := context.WithValue(request.Context(), OptionsContextKey, options)
		request = request.WithContext(ctx)

		h(response, request)
	}
}

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		header := request.Header.Get(authorizationHeader)

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			fmt.Println("headerParts", "invalid auth header")
			return
		}

		if len(headerParts[1]) == 0 {
			fmt.Println("token is empty")
			return
		}

		userId, err := ParseToken(headerParts[1])
		if err != nil {
			fmt.Println("ParseToken", err.Error())
			return
		}

		//locationId, err := strconv.Atoi(id)

		ctx := context.WithValue(request.Context(), userCtx, userId)
		request = request.WithContext(ctx)

		h(response, request)

	}
}

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

func ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}

	return claims.UserId, nil
}

// func (h *Handler) AuthMiddleware(c *gin.Context) {
// 	header := c.GetHeader(authorizationHeader)
// 	if header == "" {
// 		fmt.Println()
// 		return
// 	}

// 	headerParts := strings.Split(header, " ")
// 	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
// 		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
// 		return
// 	}

// 	if len(headerParts[1]) == 0 {
// 		newErrorResponse(c, http.StatusUnauthorized, "token is empty")
// 		return
// 	}

// 	userId, err := h.services.Authorization.ParseToken(headerParts[1])
// 	if err != nil {
// 		newErrorResponse(c, http.StatusUnauthorized, err.Error())
// 		return
// 	}

// 	c.Set(userCtx, userId)
// }
