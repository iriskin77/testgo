package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/internal/cargos"
	"github.com/iriskin77/testgo/internal/cars"
	"github.com/iriskin77/testgo/internal/files"
	"github.com/iriskin77/testgo/internal/locations"
	storage "github.com/iriskin77/testgo/pkg/db"
	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
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

	if err := initConfig(); err != nil {
		logger.Fatalf("error initialization config %s", err.Error())
	}

	logger.Info("config has been initialized")

	// Инициализруем подключение к БД

	ctx := context.Background()

	db, err := storage.NewPostgresDB(ctx, storage.ConfigDB{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logger.Fatalf("error initialization config %s", err.Error())
	}

	logger.Info("db has been initialized")

	repo := NewRepository(db, logger)
	service := NewService(repo, logger)
	handlers := NewHandler(service, logger)

	handlers.RegisterCarHandlers(s.router)
	handlers.RegisterFileHandlers(s.router)
	handlers.RegisterLocationsHandler(s.router)
	handlers.RegisterCargoHandlers(s.router)

	logger.Info("handlers have been initialized")

	logger.Info("starting API Server")

	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)

}

func initConfig() error {
	viper.AddConfigPath("/home/abc/Рабочий стол/welbexfile/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

type Repository struct {
	cars.RepositoryCar
	cargos.RepositoryCargo
	files.RepositoryFile
	locations.RepositoryLocation
}

func NewRepository(db *pgxpool.Pool, logger logging.Logger) *Repository {
	return &Repository{
		RepositoryCar:      cars.NewCarDB(db, logger),
		RepositoryCargo:    cargos.NewCargoDB(db, logger),
		RepositoryFile:     files.NewFileDB(db, logger),
		RepositoryLocation: locations.NewLocationDB(db, logger),
	}
}

type Services struct {
	cars.ServiceCar
	cargos.ServiceCargo
	files.ServiceFile
	locations.ServiceLocation
}

func NewService(repo *Repository, logger logging.Logger) *Services {
	return &Services{
		ServiceCar:      cars.NewCarService(repo.RepositoryCar, logger),
		ServiceCargo:    cargos.NewCargoService(repo.RepositoryCargo, logger),
		ServiceFile:     files.NewFileService(repo.RepositoryFile, logger),
		ServiceLocation: locations.NewLocationService(repo.RepositoryLocation, logger),
	}

}

type Handler struct {
	cars.HandlerCar
	cargos.HandlerCargo
	locations.HandlerLocation
	files.HandlerFile
}

func NewHandler(services *Services, logger logging.Logger) *Handler {
	return &Handler{
		*cars.NewHandlerCar(services.ServiceCar, logger),
		*cargos.NewHandlerCargo(services.ServiceCargo, logger),
		*locations.NewHandlerLocation(services.ServiceLocation, logger),
		*files.NewHandlerFile(services.ServiceFile, logger),
	}
}
