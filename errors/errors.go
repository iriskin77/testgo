package errors

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type ctxKey string

const (
	KeyRequestInfo ctxKey = "request_info"
)

type RequestInfo struct {
	Status  int
	Message string
}

type ClientResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewErrorClientResponse(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	response := ClientResponse{
		Status:  statusCode,
		Message: message,
	}
	sendData(ctx, w, response, statusCode, message)
}

func sendData(ctx context.Context, w http.ResponseWriter, response ClientResponse, statusCode int, message string) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	requestInfo, ok := ctx.Value(KeyRequestInfo).(*RequestInfo)
	if !ok {
		log.Println("Request info not found in context")
	} else {
		requestInfo.Status = statusCode
		requestInfo.Message = message
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
