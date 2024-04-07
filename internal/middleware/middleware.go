package middleware

import (
	"context"
	"net/http"
	"strings"
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

func Middleware(h http.HandlerFunc) http.HandlerFunc {
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
