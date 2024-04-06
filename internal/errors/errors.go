package errors

import (
	"encoding/json"
	"net/http"
)

type ClientErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewErrorClientResponse(response http.ResponseWriter, statusCode int, message string) {
	errorResponse := ClientErrorResponse{
		Status:  statusCode,
		Message: message,
	}
	SendErrorResponse(response, errorResponse)
}

func SendErrorResponse(response http.ResponseWriter, errorResponse ClientErrorResponse) {
	responseJSON, err := json.Marshal(errorResponse)
	if err != nil {
		http.Error(response, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(responseJSON)
}
