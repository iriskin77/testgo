package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) RegisterLocationsHandler(router *mux.Router) {
	router.HandleFunc("/location", h.InsertFileToDB).Methods("Post")

}

func (h *Handler) InsertFileToDB(response http.ResponseWriter, request *http.Request) {

	// vars := mux.Vars(request)
	// id := vars["id"]

	// fileId, err := strconv.Atoi(id)

	// if err != nil {
	// 	panic(err)
	// }

	h.services.Location.InsertFileToDB(2)

}
