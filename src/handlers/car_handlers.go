package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) RegisterCarHandlers(router *mux.Router) {
	router.HandleFunc("/createcar", h.CreateCar).Methods("POST")
}

func (h *Handler) CreateCar(response http.ResponseWriter, request *http.Request) {

	newCar := &models.CarRequest{}

	json.NewDecoder(request.Body).Decode(newCar)

	fmt.Println(newCar)

	car, err := h.services.Car.CreateCar(context.Background(), newCar)

	fmt.Println(car)

	if err != nil {
		logrus.Fatal("CreateCar", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	resp, err := json.Marshal(car)

	if err != nil {
		logrus.Fatal("json.Marshal(location)", err.Error())
		response.WriteHeader(http.StatusInternalServerError)
	}

	response.Write(resp)

}
