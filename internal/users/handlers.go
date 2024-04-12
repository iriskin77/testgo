package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/pkg/logging"
)

const (
	locationUrl  = "/api/user/{id}"
	locationsUrl = "/api/users"
)

type HandlerUser struct {
	services ServiceUser
	logger   logging.Logger
}

func NewHandlerUser(services ServiceUser, logger logging.Logger) *HandlerUser {
	return &HandlerUser{
		services: services,
		logger:   logger,
	}
}

func (h *HandlerUser) RegisterUserHandler(router *mux.Router) {
	router.HandleFunc(locationsUrl, h.CreateUser).Methods("Post")

}

func (h *HandlerUser) CreateUser(response http.ResponseWriter, request *http.Request) {

	newUser := &User{}

	json.NewDecoder(request.Body).Decode(newUser)

	fmt.Println("Before validation", newUser)

	if err := newUser.CreateUserValidate(); err != nil {
		h.logger.Errorf("Failed to validate: invalid data to create a new user %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	fmt.Println("After validation", newUser)

	newUserId, err := h.services.CreateUser(context.Background(), newUser)

	if err != nil {
		h.logger.Errorf("Failed to create a new user %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(newUserId)

	if err != nil {
		h.logger.Errorf("Failed to Marshal data from db in handlers %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Write(resp)

}
