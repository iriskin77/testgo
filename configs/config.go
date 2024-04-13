package configs

import "os"

type Config struct {
	Postgres ConfigPostgres

	Bindaddr string
	Loglevel string
}

type ConfigPostgres struct {
	Host     string
	Port     string
	User     string
	Password string
	NameDB   string
	SSLMode  string
}

func NewConfig() *Config {
	return &Config{
		Postgres: ConfigPostgres{
			User:     os.Getenv("POSTGRES_USER"),
			NameDB:   os.Getenv("POSTGRES_DB"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			SSLMode:  os.Getenv("SSLMODE"),
		},
		Bindaddr: os.Getenv("BINDADDR"),
		Loglevel: os.Getenv("LOGLEVEL"),
	}
}
