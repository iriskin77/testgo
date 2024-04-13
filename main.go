package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/iriskin77/testgo/configs"
	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/iriskin77/testgo/server"
	"github.com/joho/godotenv"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server Petstore server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host petstore.swagger.io
// @BasePath /v2
func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println(os.Getenv("BINDADDR"))

	// Инициализируем логгер
	logging.InitLogger()

	logger := logging.GetLogger()

	logger.Info("logger has been initialized")

	// Инициализируем конфиг
	config := configs.NewConfig()

	srv, err := server.RunServer(logger, config.Postgres, config.Bindaddr)

	if err != nil {
		logger.Fatalf("Failed to start server", err.Error())
	}

	if err = srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Errorf("HTTP server ListenAndServe", err.Error())
	}

}
