package server

type ConfigServer struct {
	BindAddr string
	LogLevel string
}

func NewConfigServer() *ConfigServer {
	return &ConfigServer{
		BindAddr: ":8000",
		LogLevel: "debug",
	}
}
