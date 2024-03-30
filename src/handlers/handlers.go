package handlers

import "github.com/iriskin77/testgo/src/services"

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}
