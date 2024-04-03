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
)

type APIServer struct {
	logger       *logrus.Logger
	serverConfig *ConfigServer
	router       *mux.Router
}

func NewApiServer(serverConfig *ConfigServer) *APIServer {
	return &APIServer{
		serverConfig: serverConfig,
		logger:       logrus.New(),
		router:       mux.NewRouter(),
	}

}

func (s *APIServer) RunServer() error {

	s.logger.Info("starting API Server")

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
		logrus.Fatal("failed to initialize db: %s", err.Error())
	}

	logrus.Info("db has been initialized")

	// repo := repository.NewRepository(db) // возвращает репозиторий (struct) Repository с методами для БД (CreateUser...)

	// service := services.NewService(repo) // возвращает сервис (struct) Service с методами для БД (CreateUser...)

	// handlers := handlers.NewHandler(service) // возвращает хэндлеры (struct) Handler

	// handlers.RegisterFileHandlers(s.router)
	// handlers.RegisterLocationsHandler(s.router)
	// handlers.RegisterCarHandlers(s.router)

	handlersCars := InitCars(db)
	handlersFiles := InitFiles(db)
	handlersLocations := InitLocations(db)
	handersCargos := InitCargo(db)

	handlersCars.RegisterCarHandlers(s.router)
	handlersFiles.RegisterFileHandlers(s.router)
	handlersLocations.RegisterLocationsHandler(s.router)
	handersCargos.RegisterCargoHandlers(s.router)

	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)

}

func initConfig() error {
	viper.AddConfigPath("/home/abc/Рабочий стол/welbexfile/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func InitCars(db *pgxpool.Pool) *cars.Handler {

	repo := cars.NewCarDB(db)
	service := cars.NewCarService(repo)
	handers := cars.NewHandler(service)

	return handers

}

func InitFiles(db *pgxpool.Pool) *files.Handler {

	repo := files.NewFileDB(db)
	service := files.NewFileService(repo)
	handers := files.NewHandler(service)

	return handers

}

func InitLocations(db *pgxpool.Pool) *locations.Handler {

	repo := locations.NewLocationDB(db)
	service := locations.NewLocationService(repo)
	handers := locations.NewHandler(service)

	return handers

}

func InitCargo(db *pgxpool.Pool) *cargos.Handler {

	repo := cargos.NewCargoDB(db)
	service := cargos.NewCargoService(repo)
	handers := cargos.NewHandler(service)

	return handers

}
