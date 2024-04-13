package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/configs"
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
	_ "github.com/lib/pq"
	httpSwagger "github.com/swaggo/http-swagger"
)

func NewHttpServer(ctx context.Context, logger logging.Logger, postgres configs.ConfigPostgres, BindAddr string) (*http.Server, error) {

	// Connecting to db

	db, err := storage.NewPostgresDB(ctx, postgres)

	if err != nil {
		logger.Fatal("Error to connect to db")
	}

	logger.Info("db has been initialized")

	// Initializing repository, services, handlers
	repo := NewRepository(db, logger)
	service := NewService(repo, logger)
	h := NewHandler(service, logger)

	router := mux.NewRouter()

	// Car handlers
	router.HandleFunc("/api/create_car", h.HandlerCar.CreateCar).Methods("POST")
	router.HandleFunc("/api/car/{id}", h.UpdateCarById).Methods("PUT")

	// File handlers
	router.HandleFunc("/api/files", h.UploadFile).Methods("POST")
	router.HandleFunc("/api/download_file/{id}", h.DownloadFile).Methods("GET")
	router.HandleFunc("/api/upload_locs_from_file/{id}", h.BulkInsertLocations).Methods("POST")

	// Location handlers
	router.HandleFunc("/api/create_location", h.CreateLocation).Methods("POST")
	router.HandleFunc("/api/get_location/{id}", h.GetLocationById).Methods("GET")
	router.HandleFunc("/api/get_locations", middleware.SortMiddleware(h.GetLocationsList)).Methods("GET")

	// Cargo handlers
	router.HandleFunc("/api/create_cargo", h.CreateCargo).Methods("POST")
	router.HandleFunc("/api/get_cargo/{id}", h.GetCargoByIDCars).Methods("GET")
	router.HandleFunc("/api/get_cargos", h.GetListCargos).Methods("GET")
	router.HandleFunc("/api/update_cargo/{id}", h.UpdateCargoById).Methods("PUT")

	// User handlers
	router.HandleFunc("/api/create_user", middleware.AuthMiddleware(h.CreateUser)).Methods("POST")
	router.HandleFunc("/api/login_user", h.LoginUser).Methods("GET")

	// swagger
	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)

	logger.Info("handlers have been initialized")

	return &http.Server{
		Addr:    BindAddr,
		Handler: router,
	}, nil
}

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
