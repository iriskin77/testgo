package locations

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/constants"
	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/internal/middleware"
	"github.com/iriskin77/testgo/pkg/logging"
)

type HandlerLocation struct {
	services ServiceLocation
	logger   logging.Logger
}

func NewHandlerLocation(services ServiceLocation, logger logging.Logger) *HandlerLocation {
	return &HandlerLocation{
		services: services,
		logger:   logger,
	}
}

// CreateLocation godoc
// @Summary Create a new location
// @Description Create a new location with the input paylod
// @Tags location
// @Accept  json
// @Produce  json
// @Param input body Location true "Create location"
// @Success 200 {integer} integer 1
// @Router /api/create_location [post]
func (h *HandlerLocation) CreateLocation(response http.ResponseWriter, request *http.Request) {

	newLocation := &Location{}

	json.NewDecoder(request.Body).Decode(newLocation)

	if err := newLocation.CreateLocationValidate(); err != nil {
		h.logger.Errorf("Failed to validate data to create location %s", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
		return
	}

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

// Get location by Id
// @Summary location id
// @Description Create a new location with the input paylod
// @Tags location
// @Accept  json
// @Produce  json
// @Param id path int true "Get location"
// @Success 200 {object} Location
// @Router /api/get_location/{id} [get]
func (h *HandlerLocation) GetLocationById(response http.ResponseWriter, request *http.Request) {

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

// GetLocationsList
// @Summary Get list
// @Description get all locations
// @Tags location
// @Accept  json
// @Produce  json
// @Success 200 {array} []Location
// @Router /api/get_locations [get]
func (h *HandlerLocation) GetLocationsList(response http.ResponseWriter, request *http.Request) {

	sortOptions := request.Context().Value(constants.OptionsContextKey).(middleware.SortOptions)

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
