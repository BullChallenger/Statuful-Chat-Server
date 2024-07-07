package service

import (
	"chat-controller/repository"
	"chat-controller/types/table"
	"fmt"
)

type Service struct {
	repository *repository.Repository
	ServerList map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	service := &Service{repository: repository, ServerList: make(map[string]bool)}
	service.setServerInfo()
	return service
}

func (service *Service) setServerInfo() {
	if serverList, err := service.ReadAvailableServerInfo(); err != nil {
		panic(err)
	} else {
		for _, server := range serverList {
			fmt.Println(server.IP)
			service.ServerList[server.IP] = true
		}
	}
}

func (service *Service) ReadAvailableServerInfo() ([]*table.ServerInfo, error) {
	return service.repository.ReadAvailableServerInfo()
}
