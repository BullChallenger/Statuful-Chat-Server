package service

import (
	"chat-controller/repository"
	"chat-controller/types/table"
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Service struct {
	repository *repository.Repository
	ServerList map[string]bool
}

func NewService(repository *repository.Repository) *Service {
	service := &Service{repository: repository, ServerList: make(map[string]bool)}
	service.setServerInfo()

	if err := service.repository.Kafka.RegisterSubTopic("chat"); err != nil {
		panic(err)
	} else {
		go service.loopSubKafka()
	}

	return service
}

func (service *Service) loopSubKafka() {
	for {
		ev := service.repository.Kafka.Poll(100)
		switch event := ev.(type) {
		case *kafka.Message:
			type ServerInfoEvent struct {
				IP     string
				Status bool
			}

			var decoder ServerInfoEvent
			if err := json.Unmarshal(event.Value, &decoder); err != nil {
				log.Println("Failed to Decode Event", event.Value)
			} else {
				fmt.Println(decoder)
				service.ServerList[decoder.IP] = decoder.Status
			}

		case *kafka.Error:
			log.Println("Failed to Polling Event", event.Error())
		}
	}
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
