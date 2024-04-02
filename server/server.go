package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iriskin77/testgo/src/handlers"
	"github.com/iriskin77/testgo/src/repository"
	"github.com/iriskin77/testgo/src/services"
	"github.com/iriskin77/testgo/storage"
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

	repo := repository.NewRepository(db) // возвращает репозиторий (struct) Repository с методами для БД (CreateUser...)

	service := services.NewService(repo) // возвращает сервис (struct) Service с методами для БД (CreateUser...)

	handlers := handlers.NewHandler(service) // возвращает хэндлеры (struct) Handler

	handlers.RegisterFileHandlers(s.router)
	handlers.RegisterLocationsHandler(s.router)

	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)

}

func initConfig() error {
	viper.AddConfigPath("/home/abc/Рабочий стол/welbexfile/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
