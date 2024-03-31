package main

import "github.com/iriskin77/testgo/server"

func main() {
	config := server.NewConfigServer()

	s := server.NewApiServer(config)

	s.RunServer()
}
