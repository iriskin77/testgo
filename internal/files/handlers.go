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
	"github.com/iriskin77/testgo/internal/errors"
	"github.com/iriskin77/testgo/pkg/logging"
)

type HandlerFile struct {
	services ServiceFile
	logger   logging.Logger
}

func NewHandlerFile(services ServiceFile, logger logging.Logger) *HandlerFile {
	return &HandlerFile{
		services: services,
		logger:   logger,
	}
}

func (h *HandlerFile) UploadFile(response http.ResponseWriter, request *http.Request) {
	request.ParseMultipartForm(10 << 20)

	// Берем файл из хэндлера
	file, handler, err := request.FormFile("file")
	if err != nil {
		h.logger.Errorf("Failed to CreateLocation in handlers %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	defer file.Close()

	// Создаем папку для хранения файлов, если ее не существует
	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		h.logger.Errorf("Failed to Create dir uploads to store files %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	// Задаем путь до файла
	absPath, _ := filepath.Abs("uploads")
	pathFile := filepath.Join(absPath, handler.Filename)
	//fmt.Println(pathFile)

	// Проверяем, есть ли в папке такой файл
	if _, err := os.Stat(pathFile); err == nil {
		h.logger.Errorf("Failed to find path to the file %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	fileExt := filepath.Ext(pathFile)

	// Проверяем формат файла
	if fileExt != ".csv" {
		h.logger.Errorf("File should be .csv %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	openedFile, err := os.OpenFile(pathFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		h.logger.Errorf("Cannot open the file %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}
	defer openedFile.Close()

	_, err = io.Copy(openedFile, file)
	if err != nil {
		h.logger.Errorf("Failed to copy file to dir uploads (file storage) %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	// Если файла с таким именем нет, то сохраняем файл в БД
	newFile := &File{
		Name:      handler.Filename,
		File_path: pathFile,
	}

	fileId, err := h.services.UploadFile(context.Background(), newFile)

	if err != nil {
		h.logger.Errorf("Failed to upload file path to DB %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	}

	if id, err := json.Marshal(fileId); err != nil {
		h.logger.Errorf("Failed to marshal file id as a response from web-server %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	} else {
		response.Write(id)
	}

}

// Download file by its id
// @Summary file id
// @Description Get a file by id
// @Tags file
// @Produce  json
// @Param id path int true "Downloading a file"
// @Router /api/download_file/{id} [get]
func (h *HandlerFile) DownloadFile(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["id"]

	fileId, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Errorf("Failed to parse file id from user request %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusNotFound, err.Error())
		return
	}

	file, err := h.services.DownloadFile(context.Background(), fileId)

	if err != nil {
		h.logger.Errorf("Failed to get file path from DB %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusNotFound, err.Error())
		return
	}

	Openfile, err := os.Open(file.File_path) //Open the file to be downloaded later

	if err != nil {
		h.logger.Errorf("Failed to get file path from DB %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusNotFound, err.Error())
		return
	}

	defer Openfile.Close() //Close after function return

	response.Header().Set("Content-Type", "application/csv")
	response.Header().Set("Content-Disposition", "attachment; filename="+file.Name)

	http.ServeFile(response, request, file.File_path)

}

// Upload list locations
// @Summary upload list locations
// @Description Uploading a list of locations from the file with zips
// @Tags file
// @Param id path int true "Uploding locations"
// @Produce  json
// @Router /api/upload_locs_from_file/{id} [post]
func (h *HandlerFile) BulkInsertLocations(response http.ResponseWriter, request *http.Request) {

	vars := mux.Vars(request)
	id := vars["id"]

	fileId, err := strconv.Atoi(id)

	if err != nil {
		h.logger.Errorf("Failed to parse file id from user request %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusNotFound, err.Error())
		return
	}

	fileIdRes, err := h.services.BulkInsertLocations(context.Background(), fileId)

	if err != nil {
		h.logger.Errorf("Failed to get file path from DB %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusNotFound, err.Error())
		return
	}

	if id, err := json.Marshal(fileIdRes); err != nil {
		h.logger.Errorf("Failed to marshal file id as a response from web-server %s", err.Error())
		errors.NewErrorClientResponse(response, http.StatusInternalServerError, err.Error())
		return
	} else {
		response.Write(id)
	}

}
