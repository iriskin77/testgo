package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/iriskin77/testgo/constants"
	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/pkg/logging"
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

// Create user
// @Summary Create a new user
// @Description Create a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param input body User true "Create user"
// @Success 200 {integer} integer 1
// @Router /api/create_user [post]
func (h *HandlerUser) CreateUser(response http.ResponseWriter, request *http.Request) {

	userIdToken, ok := request.Context().Value(constants.UserContextKey).(int)

	if !ok {
		fmt.Println(ok, "CreateUser, error")
		return
	}

	fmt.Println("userIdToken", userIdToken)

	newUser := &User{}

	json.NewDecoder(request.Body).Decode(newUser)

	if err := newUser.CreateUserValidate(); err != nil {
		h.logger.Errorf("Failed to validate: invalid data to create a new user %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

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

// login user
// @Summary login user
// @Description login user
// @Tags user
// @Accept  json
// @Produce  json
// @Param input body UserByUsernamePassword true "Create user"
// @Success 200 {integer} integer 1
// @Router /api/login_user [post]
func (h *HandlerUser) LoginUser(response http.ResponseWriter, request *http.Request) {

	userInput := &UserByUsernamePassword{}

	json.NewDecoder(request.Body).Decode(userInput)

	fmt.Println(userInput)

	if err := userInput.CreateUserSignInValidate(); err != nil {
		h.logger.Errorf("Failed to validate users input %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := h.services.GenerateToken(context.Background(), userInput.Username, userInput.Password_hash)
	if err != nil {
		h.logger.Errorf("Failed to generate token %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(token)

	resp, err := json.Marshal(token)

	if err != nil {
		h.logger.Errorf("Failed to Marshal data from db in handlers %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Write(resp)

}
