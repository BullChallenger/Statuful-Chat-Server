package app

import (
	"chat-controller/config"
	"chat-controller/network"
	"chat-controller/repository"
	"chat-controller/service"
)

type App struct {
	config *config.Config

	repository *repository.Repository
	service    *service.Service
	network    *network.Server
}

func NewApp(config *config.Config) *App {
	app := &App{config: config}

	var err error
	if app.repository, err = repository.NewRepository(config); err != nil {
		panic(err)
	} else {
		app.service = service.NewService(app.repository)
		app.network = network.NewServer(app.service, config.Info.Port)
	}
	return app
}

func (app *App) StartServer() error {
	return app.network.StartServer()
}
