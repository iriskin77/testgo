package cars

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/errors"
	"github.com/iriskin77/testgo/models"
	"go.uber.org/zap"
)

type Handler struct {
	services ServiceCar
	logger   *zap.Logger
}

func NewHandler(services ServiceCar, logger *zap.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) RegisterCarHandlers(router *mux.Router) {
	router.HandleFunc("/createcar", h.CreateCar).Methods("POST")
}

func (h *Handler) CreateCar(response http.ResponseWriter, request *http.Request) {

	newCar := &models.CarRequest{}

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
