package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/internal/cargos"
	"github.com/iriskin77/testgo/internal/cars"
	"github.com/iriskin77/testgo/internal/files"
	"github.com/iriskin77/testgo/internal/locations"
	"github.com/iriskin77/testgo/pkg/logging"
	"github.com/iriskin77/testgo/storage"
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

	handlersCars := InitCars(db, logger)
	handlersFiles := InitFiles(db, logger)
	handlersLocations := InitLocations(db, logger)
	handersCargos := InitCargo(db, logger)

	handlersCars.RegisterCarHandlers(s.router)
	handlersFiles.RegisterFileHandlers(s.router)
	handlersLocations.RegisterLocationsHandler(s.router)
	handersCargos.RegisterCargoHandlers(s.router)

	logger.Info("handlers have been initialized")

	logger.Info("starting API Server")

	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)

}

func initConfig() error {
	viper.AddConfigPath("/home/abc/Рабочий стол/welbexfile/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func InitCars(db *pgxpool.Pool, logger logging.Logger) *cars.Handler {

	repo := cars.NewCarDB(db, logger)
	service := cars.NewCarService(repo, logger)
	handers := cars.NewHandler(service, logger)

	return handers

}

func InitFiles(db *pgxpool.Pool, logger logging.Logger) *files.Handler {

	repo := files.NewFileDB(db, logger)
	service := files.NewFileService(repo, logger)
	handers := files.NewHandler(service, logger)

	return handers

}

func InitLocations(db *pgxpool.Pool, logger logging.Logger) *locations.Handler {

	repo := locations.NewLocationDB(db, logger)
	service := locations.NewLocationService(repo, logger)
	handers := locations.NewHandler(service, logger)

	return handers

}

func InitCargo(db *pgxpool.Pool, logger logging.Logger) *cargos.Handler {

	repo := cargos.NewCargoDB(db, logger)
	service := cargos.NewCargoService(repo, logger)
	handers := cargos.NewHandler(service, logger)

	return handers

}
