package server

import (
	"net/http"

	"github.com/gorilla/mux"
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

	_, err := storage.NewPostgresDB(storage.ConfigDB{
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

	return http.ListenAndServe(s.serverConfig.BindAddr, s.router)

}

func initConfig() error {
	//path := filepath.Base("config.yml")
	//abs_path, _ := filepath.Abs(path)
	//p := "/home/abc/Рабочий стол/welbexfile/configs"
	//abs_path := filepath.Base("configs")
	viper.AddConfigPath("/home/abc/Рабочий стол/welbexfile/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
