package files

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/errors"
	"github.com/iriskin77/testgo/models"
	"go.uber.org/zap"
)

const (
	file         = "/files"
	filedownload = "/files/{id}"
)

type Handler struct {
	services ServiceFile
	logger   *zap.Logger
}

func NewHandler(services ServiceFile, logger *zap.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) RegisterFileHandlers(router *mux.Router) {
	router.HandleFunc(file, h.UploadFile).Methods("POST")
	router.HandleFunc(filedownload, h.DownloadFile).Methods("GET")
}

func (h *Handler) UploadFile(response http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)

	// Берем файл из хэндлера
	file, handler, err := request.FormFile("file")
	if err != nil {
		h.logger.Error("Failed to CreateLocation in handlers", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	defer file.Close()

	// Создаем папку для хранения файлов, если ее не существует
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		h.logger.Error("Failed to Create dir uploads to store files", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	// Задаем путь до файла
	absPath, _ := filepath.Abs("uploads")
	pathFile := filepath.Join(absPath, handler.Filename)
	//fmt.Println(pathFile)

	// Проверяем, есть ли в папке такой файл
	if _, err := os.Stat(pathFile); err == nil {
		h.logger.Error("", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	fileExt := filepath.Ext(pathFile)

	// Проверяем формат файла
	if fileExt != ".csv" {
		h.logger.Error("File should be .csv", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	openedFile, err := os.OpenFile(pathFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		h.logger.Error("Cannot open the file", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}
	defer openedFile.Close()

	_, err = io.Copy(openedFile, file)
	if err != nil {
		h.logger.Error("Failed to copy file to dir uploads (file storage)", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	// Если файла с таким именем нет, то сохраняем файл в БД
	newFile := &models.File{
		Name:      handler.Filename,
		File_path: pathFile,
	}

	fileId, err := h.services.UploadFile(context.Background(), newFile)

	if err != nil {
		h.logger.Error("Failed to upload file path to DB", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
		return
	}

	if id, err := json.Marshal(fileId); err != nil {
		h.logger.Error("Failed to marshal file id as a response from web-server", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusInternalServerError, err.Error())
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
		h.logger.Error("Failed to parse file id from user request", zap.Error(err))
		errors.NewErrorClientResponse(request.Context(), response, http.StatusNotFound, err.Error())
		return
	}

	file, err := h.services.DownloadFile(context.Background(), fileId)

	if err != nil {
		h.logger.Error("Failed to get file path from DB", zap.Error(err)) //return 404 if file is not found
		errors.NewErrorClientResponse(request.Context(), response, http.StatusNotFound, err.Error())
		return
	}

	Openfile, err := os.Open(file.File_path) //Open the file to be downloaded later

	if err != nil {
		h.logger.Error("Failed to get file path from DB", zap.Error(err)) //return 404 if file is not found
		errors.NewErrorClientResponse(request.Context(), response, http.StatusNotFound, err.Error())
		return
	}

	defer Openfile.Close() //Close after function return

	response.Header().Set("Content-Type", "application/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	http.ServeFile(response, request, file.File_path)

}
