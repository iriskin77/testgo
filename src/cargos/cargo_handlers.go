package cargos

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/models"
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
}

func (h *Handler) CreateCargo(response http.ResponseWriter, request *http.Request) {

	newCargo := &models.CargoRequest{}

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
