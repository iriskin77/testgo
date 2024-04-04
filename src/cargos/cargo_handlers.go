package cargos

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/errors"
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

func (h *Handler) RegisterCargoHandlers(router *mux.Router) {
	router.HandleFunc("/createcargo", h.CreateCargo).Methods("POST")
	router.HandleFunc("/get_cargo_cars/{id}", h.GetCargoByIDCars).Methods("GET")
	router.HandleFunc("/get_list_cargos", h.GetListCargos).Methods("GET")
}

func (h *Handler) CreateCargo(response http.ResponseWriter, request *http.Request) {

	newCargo := &CargoRequest{}

	json.NewDecoder(request.Body).Decode(newCargo)

	cargoId, err := h.services.CreateCargo(context.Background(), newCargo)

	if err != nil {
		h.logger.Error("Failed to CreateCargo", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(cargoId)

	if err != nil {
		h.logger.Error("Failed to Marshal response (cargo id)", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}

func (h *Handler) GetCargoByIDCars(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["id"]

	cargoId, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Error("Failed to parse id from client request", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.services.GetCargoCars(context.Background(), cargoId)

	if err != nil {
		h.logger.Error("Failed to get the cargo and the closest cars from DB", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(res)

	if err != nil {
		h.logger.Error("Failed to parse response (cargo and cars)", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)
}

func (h *Handler) GetListCargos(response http.ResponseWriter, request *http.Request) {

	res, err := h.services.GetListCargos(context.Background())

	if err != nil {
		h.logger.Error("Failed to get list cargos and the closest cars from DB", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(res)

	if err != nil {
		h.logger.Error("Failed to parse response (cargo and cars)", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}
