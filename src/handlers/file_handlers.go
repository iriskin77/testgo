package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/models"
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
	request.ParseMultipartForm(10 << 20)

	// Берем файл из хэндлера
	file, handler, err := request.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()

	// Создаем папку для хранения файлов, если ее не существует
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// Задаем путь до файла
	absPath, _ := filepath.Abs("uploads")
	pathFile := filepath.Join(absPath, handler.Filename)
	//fmt.Println(pathFile)

	// Проверяем, есть ли в папке такой файл
	if _, err := os.Stat(pathFile); err == nil {
		http.Error(response, "File already exists. You should change filename", http.StatusInternalServerError)
		return
	}

	// Если файла с таким именем нет, то сохраняем файл в БД
	newFile := &models.File{
		Name: handler.Filename,
		File: pathFile,
	}

	fileId := h.services.File.UploadFile(newFile)

	if id, err := json.Marshal(fileId); err != nil {
		http.Error(response, "Wrong", http.StatusInternalServerError)
		return
	} else {
		response.Write(id)
	}

}

func (h *Handler) DownloadFile(response http.ResponseWriter, request *http.Request) {

}
