package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/src/cargos"
	"github.com/iriskin77/testgo/src/cars"
	"github.com/iriskin77/testgo/src/files"
	"github.com/iriskin77/testgo/src/locations"
	"github.com/iriskin77/testgo/storage"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	logger, err := initLogger()
	if err != nil {
		return err
	}

	logger.Info("starting API Server")

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initialization config %s", err.Error())
	}

	// fmt.Println(viper.GetString("db.host"))

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
		logger.Error("failed to initialize db: %s", zap.Error(err))
	}

	logger.Info("db has been initialized")

	handlersCars := InitCars(db, logger)
	handlersFiles := InitFiles(db, logger)
	handlersLocations := InitLocations(db, logger)
	handersCargos := InitCargo(db)

	handlersCars.RegisterCarHandlers(s.router)
	handlersFiles.RegisterFileHandlers(s.router)
	handlersLocations.RegisterLocationsHandler(s.router)
	handersCargos.RegisterCargoHandlers(s.router)

	logger.Info("handlers have been registered")

	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)

}

func initLogger() (*zap.Logger, error) {
	return zap.Config{
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build()
}

func initConfig() error {
	viper.AddConfigPath("/home/abc/Рабочий стол/welbexfile/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func InitCars(db *pgxpool.Pool, logger *zap.Logger) *cars.Handler {

	repo := cars.NewCarDB(db, logger)
	service := cars.NewCarService(repo, logger)
	handers := cars.NewHandler(service, logger)

	return handers

}

func InitFiles(db *pgxpool.Pool, logger *zap.Logger) *files.Handler {

	repo := files.NewFileDB(db, logger)
	service := files.NewFileService(repo, logger)
	handers := files.NewHandler(service, logger)

	return handers

}

func InitLocations(db *pgxpool.Pool, logger *zap.Logger) *locations.Handler {

	repo := locations.NewLocationDB(db, logger)
	service := locations.NewLocationService(repo, logger)
	handers := locations.NewHandler(service, logger)

	return handers

}

func InitCargo(db *pgxpool.Pool) *cargos.Handler {

	repo := cargos.NewCargoDB(db)
	service := cargos.NewCargoService(repo)
	handers := cargos.NewHandler(service)

	return handers

}
