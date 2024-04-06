package errors

import (
	"encoding/json"
	"net/http"
)

type ClientResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewErrorClientResponse(w http.ResponseWriter, statusCode int, message string) {
	response := ClientResponse{
		Status:  statusCode,
		Message: message,
	}
	sendData(w, response)
}

func sendData(w http.ResponseWriter, response ClientResponse) {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
