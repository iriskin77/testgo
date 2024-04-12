package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/iriskin77/testgo/constants"
	"github.com/iriskin77/testgo/internal/users"
)

// Middleware for sorting

type SortOptions struct {
	Field string
	Order string
}

func SortMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		sortBy := request.URL.Query().Get("sort_by")
		sortOrder := request.URL.Query().Get("sort_order")

		if sortBy == "" {
			sortBy = constants.DefaulSortField
		}
		if sortOrder == "" {
			sortOrder = constants.DefaultSortOrder
		} else {
			upperSortOrder := strings.ToUpper(sortOrder)
			if upperSortOrder != constants.AscSort && upperSortOrder != constants.DescSort {
				response.WriteHeader(http.StatusBadRequest)
				response.Write([]byte("Нельзя так делать"))

			}
		}

		options := SortOptions{
			Field: sortBy,
			Order: sortOrder,
		}

		ctx := context.WithValue(request.Context(), constants.OptionsContextKey, options)
		request = request.WithContext(ctx)

		h(response, request)
	}
}

// Middleware to auth user. It gets user id from token
func AuthMiddleware(h http.HandlerFunc) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		header := request.Header.Get(constants.AuthorizationHeader)

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			fmt.Println("headerParts", "invalid auth header")
			return
		}

		if len(headerParts[1]) == 0 {
			fmt.Println("token is empty")
			return
		}

		userId, err := users.GetUserFromToken(headerParts[1])
		if err != nil {
			fmt.Println("ParseToken", err.Error())
			return
		}

		fmt.Println(AuthMiddleware, userId)

		ctx := context.WithValue(request.Context(), constants.UserContextKey, userId)
		request = request.WithContext(ctx)

		h(response, request)

	}
}
