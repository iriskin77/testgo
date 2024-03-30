package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	file         = "/file"
	downloadFile = "file/:id"
)

func (h *Handler) RegisterFileHandlers(router *mux.Router) {
	router.HandleFunc(file, h.UploadFile).Methods("POST")
	router.HandleFunc(downloadFile, h.DownloadFile).Methods("GET")

}

func (h *Handler) UploadFile(response http.ResponseWriter, request *http.Request) {
	fmt.Println(request.MultipartForm.File)

}

func (h *Handler) DownloadFile(response http.ResponseWriter, request *http.Request) {

}
