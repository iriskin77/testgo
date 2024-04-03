package cargos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	services ServiceCar
}

func NewHandler(services ServiceCar) *Handler {
	return &Handler{services: services}
}

func (h *Handler) RegisterCargoHandlers(router *mux.Router) {
	router.HandleFunc("/createcargo", h.CreateCargo).Methods("POST")
	router.HandleFunc("/get_cargo_cars/{id}", h.GetCargoCars).Methods("GET")
}

func (h *Handler) CreateCargo(response http.ResponseWriter, request *http.Request) {

	newCargo := &CargoRequest{}

	json.NewDecoder(request.Body).Decode(newCargo)

	cargoId, err := h.services.CreateCargo(context.Background(), newCargo)

	fmt.Println(cargoId)

	if err != nil {
		logrus.Fatal("CreateCargo", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(cargoId)

	if err != nil {
		logrus.Fatal("json.Marshal(location)", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.Write(resp)

}

func (h *Handler) GetCargoCars(response http.ResponseWriter, request *http.Request) {

	//cargo := &models.CargoCarsResponse{}

	vars := mux.Vars(request)
	id := vars["id"]

	cargoId, err := strconv.Atoi(id)

	if err != nil {
		logrus.Fatal("(h *Handler) GetLocationById", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	res, err := h.services.GetCargoCars(context.Background(), cargoId)

	if err != nil {
		fmt.Println(err.Error())
	}

	resp, err := json.Marshal(res)

	if err != nil {
		logrus.Fatal("json.Marshal(location)", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.Write(resp)

}
