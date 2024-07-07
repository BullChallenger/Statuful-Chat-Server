package main

import (
	"chat-controller/cmd/app"
	"chat-controller/config"
	"flag"
)

var pathFlag = flag.String("config", "./config.toml", "config set")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)

	a := app.NewApp(c)
	a.StartServer()
}
