package cargos

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/pkg/logging"
)

const (
	cargoUrl  = "/api/cargo/{id}"
	cargosUrl = "/api/cargos"
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

func (h *Handler) RegisterCargoHandlers(router *mux.Router) {
	router.HandleFunc(cargosUrl, h.CreateCargo).Methods("POST")
	router.HandleFunc(cargoUrl, h.GetCargoByIDCars).Methods("GET")
	router.HandleFunc(cargosUrl, h.GetListCargos).Methods("GET")
}

func (h *Handler) CreateCargo(response http.ResponseWriter, request *http.Request) {

	newCargo := &CargoRequest{}
	defer request.Body.Close()

	if err := json.NewDecoder(request.Body).Decode(newCargo); err != nil {
		h.logger.Errorf("Failed to decode request data(CreateCargo) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
	}

	cargoId, err := h.services.CreateCargo(context.Background(), newCargo)
	if err != nil {
		h.logger.Errorf("Failed to CreateCargo %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(cargoId)
	if err != nil {
		h.logger.Errorf("Failed to Marshal response (cargo id) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)
}

func (h *Handler) GetCargoByIDCars(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["id"]

	cargoId, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Errorf("Failed to parse id from client request %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	res, err := h.services.GetCargoCars(context.Background(), cargoId)

	if err != nil {
		h.logger.Errorf("Failed to get the cargo and the closest cars from DB %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(res)

	if err != nil {
		h.logger.Errorf("Failed to parse response (cargo and cars) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)
}

func (h *Handler) GetListCargos(response http.ResponseWriter, request *http.Request) {

	listCargos, err := h.services.GetListCargos(context.Background())

	if err != nil {
		h.logger.Errorf("Failed to get list cargos and the closest cars from DB %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(listCargos)

	if err != nil {
		h.logger.Errorf("Failed to parse response (cargo and cars) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}

func (h *Handler) UpdateCargoById(response http.ResponseWriter, request *http.Request) {

	cargoUpdated := &CargoUpdateRequest{}

	json.NewDecoder(request.Body).Decode(cargoUpdated)

	carUpdatedId, err := h.services.UpdateCargoById(context.Background(), cargoUpdated)

	if err != nil {
		h.logger.Errorf("Failed to update cargo %s", err.Error())
	}

	resp, err := json.Marshal(carUpdatedId)

	if err != nil {
		h.logger.Errorf("Failed to parse response (cargo and cars) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}
