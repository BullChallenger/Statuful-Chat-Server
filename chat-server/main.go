package main

import (
	"chat-server/config"
	"chat-server/network"
	"chat-server/repository"
	"chat-server/service"
	"flag"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", "localhost:8080", "port set")

func main() {
	flag.Parse()

	c := config.NewConfig(*pathFlag)
	if repository, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		n := network.NewServer(service.NewService(repository), *port)
		n.StartServer()
	}
}
