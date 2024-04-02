package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) RegisterLocationsHandler(router *mux.Router) {
	router.HandleFunc("/location", h.CreateLocation).Methods("Post")
	router.HandleFunc("/location/{id}", h.GetLocationById).Methods("Get")
	router.HandleFunc("/locations", h.GetLocationsList).Methods("Get")
}

func (h *Handler) CreateLocation(response http.ResponseWriter, request *http.Request) {

	newLocation := &models.Location{}

	json.NewDecoder(request.Body).Decode(newLocation)

	location, err := h.services.Location.CreateLocation(context.Background(), newLocation)

	if err != nil {
		logrus.Fatal("h.services.CreateLocation(newLocation)")
		response.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(location)

	if err != nil {
		logrus.Fatal("json.Marshal(location)")
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.Write(resp)

}

func (h *Handler) GetLocationById(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["id"]

	locationId, err := strconv.Atoi(id)

	if err != nil {
		logrus.Fatal("(h *Handler) GetLocationById", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Println(locationId)

	locationById, err := h.services.Location.GetLocationById(context.Background(), locationId)

	if err != nil {
		logrus.Fatal("h.services.CreateLocation(newLocation)", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(locationById)

	if err != nil {
		logrus.Fatal("json.Marshal(location)", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.Write(resp)

}

func (h *Handler) GetLocationsList(response http.ResponseWriter, request *http.Request) {

	locationsList, err := h.services.Location.GetLocationsList(context.Background())

	if err != nil {
		logrus.Fatal("(h *Handler) GetLocationsList", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(locationsList)

	if err != nil {
		logrus.Fatal("json.Marshal(location)", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.Write(resp)

}
