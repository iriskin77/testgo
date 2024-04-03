package cargos

import (
	"net/http"

	"github.com/gorilla/mux"
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

}
