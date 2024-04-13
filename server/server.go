package server

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/iriskin77/testgo/docs"
	"github.com/iriskin77/testgo/internal/cargos"
	"github.com/iriskin77/testgo/internal/cars"
	"github.com/iriskin77/testgo/internal/files"
	"github.com/iriskin77/testgo/internal/locations"
	"github.com/iriskin77/testgo/internal/middleware"
	"github.com/iriskin77/testgo/internal/users"
	storage "github.com/iriskin77/testgo/pkg/db"
	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type APIServer struct {
	serverConfig *ConfigServer
	router       *mux.Router
}

func NewApiServer(serverConfig *ConfigServer) *APIServer {
	return &APIServer{
		serverConfig: serverConfig,
		router:       mux.NewRouter(),
	}

}

func (s *APIServer) RunServer() error {

	// Инициализируем логгер
	logging.InitLogger()

	logger := logging.GetLogger()

	logger.Info("logger has been initialized")

	// if err := initConfig(); err != nil {
	// 	logger.Fatalf("error initialization config %s", err.Error())
	// }

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	logger.Info("config has been initialized")

	// Инициализруем подключение к БД

	ctx := context.Background()

	db, err := storage.NewPostgresDB(ctx, storage.ConfigDB{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Username: os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("SSLMODE"),
	})

	if err != nil {
		logger.Fatalf("error initialization config %s", err.Error())
	}

	logger.Info("db has been initialized")

	repo := NewRepository(db, logger)
	service := NewService(repo, logger)
	h := NewHandler(service, logger)

	// Car handlers
	s.router.HandleFunc("/api/cars", h.HandlerCar.CreateCar).Methods("POST")
	s.router.HandleFunc("/api/car/{id}", h.UpdateCarById).Methods("PUT")

	// File handlers
	s.router.HandleFunc("/api/files", h.UploadFile).Methods("POST")
	s.router.HandleFunc("/api/file/{id}", h.DownloadFile).Methods("GET")
	s.router.HandleFunc("/api/upload_file/{id}", h.BulkInsertLocations).Methods("PUT")

	// Location handlers
	s.router.HandleFunc("/api/createlocation", h.CreateLocation).Methods("POST")
	s.router.HandleFunc("/api/get_location/{id}", h.GetLocationById).Methods("GET")
	s.router.HandleFunc("/api/get_locations", middleware.SortMiddleware(h.GetLocationsList)).Methods("GET")

	s.router.HandleFunc("/api/createcargo", h.CreateCargo).Methods("POST")
	s.router.HandleFunc("/api/get_cargo/{id}", h.GetCargoByIDCars).Methods("GET")
	s.router.HandleFunc("/api/get_cargos", h.GetListCargos).Methods("GET")

	// User handlers
	s.router.HandleFunc("/api/create_user", middleware.AuthMiddleware(h.CreateUser)).Methods("POST")
	s.router.HandleFunc("/api/login_user", h.LoginUser).Methods("GET")

	s.router.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	)).Methods("GET")

	logger.Info("handlers have been initialized")

	logger.Info("starting API Server")

	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)

}

// func initConfig() error {
// 	pathConfig, _ := filepath.Abs("configs")
// 	fmt.Println(pathConfig)
// 	viper.AddConfigPath(pathConfig)
// 	viper.SetConfigName("config")
// 	return viper.ReadInConfig()
// }

type Repository struct {
	cars.RepositoryCar
	cargos.RepositoryCargo
	files.RepositoryFile
	locations.RepositoryLocation
	users.RepositoryUser
}

func NewRepository(db *pgxpool.Pool, logger logging.Logger) *Repository {
	return &Repository{
		RepositoryCar:      cars.NewCarDB(db, logger),
		RepositoryCargo:    cargos.NewCargoDB(db, logger),
		RepositoryFile:     files.NewFileDB(db, logger),
		RepositoryLocation: locations.NewLocationDB(db, logger),
		RepositoryUser:     users.NewUserDB(db, logger),
	}
}

type Services struct {
	cars.ServiceCar
	cargos.ServiceCargo
	files.ServiceFile
	locations.ServiceLocation
	users.ServiceUser
}

func NewService(repo *Repository, logger logging.Logger) *Services {
	return &Services{
		ServiceCar:      cars.NewCarService(repo.RepositoryCar, logger),
		ServiceCargo:    cargos.NewCargoService(repo.RepositoryCargo, logger),
		ServiceFile:     files.NewFileService(repo.RepositoryFile, logger),
		ServiceLocation: locations.NewLocationService(repo.RepositoryLocation, logger),
		ServiceUser:     users.NewUserService(repo.RepositoryUser, logger),
	}

}

type Handler struct {
	cars.HandlerCar
	cargos.HandlerCargo
	locations.HandlerLocation
	files.HandlerFile
	users.HandlerUser
}

func NewHandler(services *Services, logger logging.Logger) *Handler {
	return &Handler{
		*cars.NewHandlerCar(services.ServiceCar, logger),
		*cargos.NewHandlerCargo(services.ServiceCargo, logger),
		*locations.NewHandlerLocation(services.ServiceLocation, logger),
		*files.NewHandlerFile(services.ServiceFile, logger),
		*users.NewHandlerUser(services.ServiceUser, logger),
	}
}
