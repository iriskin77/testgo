package cars

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/pkg/logging"
	"go.uber.org/zap"
)

type Handler struct {
	services ServiceCar
	logger   logging.Logger
}

func NewHandler(services ServiceCar, logger logging.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) RegisterCarHandlers(router *mux.Router) {
	router.HandleFunc("/createcar", h.CreateCar).Methods("POST")
	router.HandleFunc("/update_car", h.UpdateCarById).Methods("PUT")
}

func (h *Handler) CreateCar(response http.ResponseWriter, request *http.Request) {

	newCar := &CarRequest{}

	json.NewDecoder(request.Body).Decode(newCar)

	carId, err := h.services.CreateCar(context.Background(), newCar)

	if err != nil {
		h.logger.Error("Failed to CreateCar", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(carId)

	if err != nil {
		h.logger.Error("Failed to Marshal response (car id)", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}

func (h *Handler) UpdateCarById(response http.ResponseWriter, request *http.Request) {

	carUpdatedData := &CarUpdateRequest{}

	json.NewDecoder(request.Body).Decode(carUpdatedData)

	carUpdatedId, err := h.services.UpdateCarById(context.Background(), carUpdatedData)

	if err != nil {
		h.logger.Error("Failed to update car by id", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusNotFound, err.Error())
		return
	}

	resp, err := json.Marshal(carUpdatedId)

	if err != nil {
		h.logger.Error("Failed to Marshal response (car id)", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}
