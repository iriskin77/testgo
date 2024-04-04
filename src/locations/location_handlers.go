package locations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/models"
	"github.com/iriskin77/testgo/src/errors"
	"go.uber.org/zap"
)

type Handler struct {
	services ServiceLocation
	logger   *zap.Logger
}

func NewHandler(services ServiceLocation, logger *zap.Logger) *Handler {
	return &Handler{services: services, logger: logger}
}

func (h *Handler) RegisterLocationsHandler(router *mux.Router) {
	router.HandleFunc("/location", h.CreateLocation).Methods("Post")
	router.HandleFunc("/location/{id}", h.GetLocationById).Methods("Get")
	router.HandleFunc("/locations", h.GetLocationsList).Methods("Get")
}

func (h *Handler) CreateLocation(response http.ResponseWriter, request *http.Request) {

	newLocation := &models.Location{}

	json.NewDecoder(request.Body).Decode(newLocation)

	location, err := h.services.CreateLocation(context.Background(), newLocation)

	if err != nil {
		h.logger.Error("Failed to CreateLocation in handlers", zap.Error(err))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(location)

	if err != nil {
		h.logger.Error("Failed to Marshal data from db in handlers", zap.Error(err))
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	response.Write(resp)

}

func (h *Handler) GetLocationById(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["id"]

	locationId, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Error("Failed to retrieve GetLocationById from db in handlers", zap.Error(err))
		http.Error(response, err.Error(), http.StatusInternalServerError)

	}

	fmt.Println(locationId)

	locationById, err := h.services.GetLocationById(context.Background(), locationId)

	if err != nil {
		h.logger.Error("Failed to retrieve GetLocationById from db in handlers", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(locationById)

	if err != nil {
		h.logger.Error("Failed to marshal GetLocationById from db in handlers", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}

func (h *Handler) GetLocationsList(response http.ResponseWriter, request *http.Request) {

	locationsList, err := h.services.GetLocationsList(context.Background())

	if err != nil {
		h.logger.Error("Failed to retrieve GetLocationsList from db in handlers", zap.Error(err))
		response.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(locationsList)

	if err != nil {
		h.logger.Error("Failed to marshal GetLocationsList from db in handlers", zap.Error(err))
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.Write(resp)

}
