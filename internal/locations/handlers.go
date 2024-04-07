package locations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/internal/middleware"
	"github.com/iriskin77/testgo/pkg/logging"
)

const (
	locationUrl  = "/api/location/{id}"
	locationsUrl = "/api/locations"
)

type Handler struct {
	services ServiceLocation
	logger   logging.Logger
}

func NewHandler(services ServiceLocation, logger logging.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) RegisterLocationsHandler(router *mux.Router) {
	router.HandleFunc(locationsUrl, h.CreateLocation).Methods("Post")
	router.HandleFunc(locationUrl, h.GetLocationById).Methods("Get")
	router.HandleFunc(locationsUrl, middleware.Middleware(h.GetLocationsList)).Methods("Get")
}

func (h *Handler) CreateLocation(response http.ResponseWriter, request *http.Request) {

	newLocation := &Location{}

	json.NewDecoder(request.Body).Decode(newLocation)

	location, err := h.services.CreateLocation(context.Background(), newLocation)

	if err != nil {
		h.logger.Errorf("Failed to CreateLocation in handlers %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(location)

	if err != nil {
		h.logger.Errorf("Failed to Marshal data from db in handlers %s", err.Error())
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
		h.logger.Errorf("Failed to retrieve GetLocationById from db in handlers %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())

	}

	fmt.Println(locationId)

	locationById, err := h.services.GetLocationById(context.Background(), locationId)

	if err != nil {
		h.logger.Errorf("Failed to retrieve GetLocationById from db in handlers %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	resp, err := json.Marshal(locationById)

	if err != nil {
		h.logger.Errorf("Failed to marshal GetLocationById from db in handlers %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	response.Write(resp)

}

func (h *Handler) GetLocationsList(response http.ResponseWriter, request *http.Request) {

	sortOptions := request.Context().Value(middleware.OptionsContextKey).(middleware.SortOptions)

	fmt.Println("sortOptions", sortOptions)

	locationsList, err := h.services.GetLocationsList(context.Background(), sortOptions)

	if err != nil {
		h.logger.Errorf("Failed to retrieve GetLocationsList from db in handlers %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
	}

	resp, err := json.Marshal(locationsList)

	if err != nil {
		h.logger.Errorf("Failed to marshal GetLocationsList from db in handlers %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
	}

	response.Header().Set("Content-Type", "application/json")
	response.Write(resp)

}
