package service

import (
	"chat-server/repository"
	"chat-server/types/schema"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

type Service struct {
	repository *repository.Repository
}

func NewService(repository *repository.Repository) *Service {
	service := &Service{repository: repository}
	return service
}

func (service *Service) ServerSet(ip string, available bool) error {
	if err := service.repository.ServerSet(ip, available); err != nil {
		fmt.Println("Failed To Server Set", "ip", ip, "available", available)
		return err
	} else {
		return nil
	}
}

func (service *Service) PublishEvent(topic string, value []byte, ch chan kafka.Event) (kafka.Event, error) {
	return service.repository.Kafka.PublishEvent(topic, value, ch)
}

func (service *Service) InsertChatting(user, message, roomName string) {
	if err := service.repository.InsertChatting(user, message, roomName); err != nil {
		log.Println("Failed to Insert Chatting", "err", err.Error())
	}
}

// EnterRoom
func (service *Service) EnterRoom(roomName string) ([]*schema.Chat, error) {
	if res, err := service.repository.ReadChatList(roomName); err != nil {
		log.Println("Failed to Get All Chat List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

// RoomList
func (service *Service) RoomList() ([]*schema.Room, error) {
	if res, err := service.repository.RoomList(); err != nil {
		log.Println("Failed to Get All Room List", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}

// MakeRoom
func (service *Service) MakeRoom(roomName string) error {
	if err := service.repository.MakeRoom(roomName); err != nil {
		log.Println("Failed to Make New Room", "err", err.Error())
		return err
	} else {
		return nil
	}
}

// Room
func (service *Service) Room(roomName string) (*schema.Room, error) {
	if res, err := service.repository.Room(roomName); err != nil {
		log.Println("Failed to Get Room", "err", err.Error())
		return nil, err
	} else {
		return res, nil
	}
}
