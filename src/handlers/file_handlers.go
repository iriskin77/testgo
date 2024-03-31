package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/models"
)

const (
	file         = "/files"
	filedownload = "/files/{id}"
)

func (h *Handler) RegisterFileHandlers(router *mux.Router) {
	router.HandleFunc(file, h.UploadFile).Methods("POST")
	router.HandleFunc(filedownload, h.DownloadFile).Methods("GET")

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

	fileExt := filepath.Ext(pathFile)

	// Проверяем формат файла
	if fileExt != ".csv" {
		http.Error(response, "File should be csv", http.StatusForbidden)
		return
	}

	openedFile, err := os.OpenFile(pathFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode("something went wrong(OpenedFile)")
		return
	}
	defer openedFile.Close()

	_, err = io.Copy(openedFile, file)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(response).Encode("something went wrong (Copy)")
		return
	}

	// Если файла с таким именем нет, то сохраняем файл в БД
	newFile := &models.File{
		Name:      handler.Filename,
		File_path: pathFile,
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

	vars := mux.Vars(request)
	id := vars["id"]

	fileId, err := strconv.Atoi(id)

	if err != nil {
		panic(err)
	}

	file, err := h.services.File.DownloadFile(fileId)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError) //return 404 if file is not found
		return
	}

	Openfile, err := os.Open(file.File_path) //Open the file to be downloaded later

	if err != nil {
		http.Error(response, err.Error(), http.StatusNotFound) //return 404 if file is not found
		return
	}

	defer Openfile.Close() //Close after function return

	response.Header().Set("Content-Type", "application/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	http.ServeFile(response, request, file.File_path)

}
