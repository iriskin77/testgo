package cars

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/pkg/logging"
)

type HandlerCar struct {
	services ServiceCar
	logger   logging.Logger
}

func NewHandlerCar(services ServiceCar, logger logging.Logger) *HandlerCar {
	return &HandlerCar{
		services: services,
		logger:   logger,
	}
}

// CreateCar
// @Summary create a new car
// @Description Create a new car with the input paylod
// @Tags car
// @Accept  json
// @Produce  json
// @Param input body CarCreateRequest true "Create car"
// @Success 200 {integer} integer 1
// @Router /api/create_car [post]
func (h *HandlerCar) CreateCar(response http.ResponseWriter, request *http.Request) {

	newCar := &CarCreateRequest{}

	json.NewDecoder(request.Body).Decode(newCar)

	carId, err := h.services.CreateCar(context.Background(), newCar)

	if err != nil {
		h.logger.Errorf("Failed to CreateCar %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(carId)

	if err != nil {
		h.logger.Errorf("Failed to Marshal response (car id) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}

// UpdateCar godoc
// @Summary Update a car
// @Description Update a car with the input paylod
// @Tags car
// @Accept  json
// @Produce  json
// @Param input body CarUpdateRequest true "Update car"
// @Success 200 {integer} integer 1
// @Router /api/create_car [put]
func (h *HandlerCar) UpdateCarById(response http.ResponseWriter, request *http.Request) {

	carUpdatedData := &CarUpdateRequest{}

	json.NewDecoder(request.Body).Decode(carUpdatedData)

	carUpdatedId, err := h.services.UpdateCarById(context.Background(), carUpdatedData)

	if err != nil {
		h.logger.Errorf("Failed to update car by id %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusNotFound, err.Error())
		return
	}

	resp, err := json.Marshal(carUpdatedId)

	if err != nil {
		h.logger.Errorf("Failed to Marshal response (car id) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}
