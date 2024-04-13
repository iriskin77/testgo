package main

import (
	"fmt"
	"log"
	"os"

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

	s3Bucket := os.Getenv("POSTGRES_USER")
	secretKey := os.Getenv("POSTGRES_DB")

	fmt.Println(s3Bucket)
	fmt.Println(secretKey)

	config := server.NewConfigServer()

	s := server.NewApiServer(config)

	s.RunServer()
}
