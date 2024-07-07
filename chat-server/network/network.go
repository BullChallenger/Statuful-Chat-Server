package network

import (
	"chat-server/service"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	engine  *gin.Engine
	service *service.Service
	port    string
	ip      string
}

func NewServer(service *service.Service, port string) *Server {
	s := &Server{engine: gin.New(), service: service, port: port}
	s.engine.Use(gin.Logger())
	s.engine.Use(gin.Recovery())
	s.engine.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	registerServer(s)
	return s
}

func (s *Server) setServerInfo() {
	if addresses, err := net.InterfaceAddrs(); err != nil {
		panic(err)
	} else {
		var ip net.IP

		for _, addr := range addresses {
			if ipnet, ok := addr.(*net.IPNet); ok {
				if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ip = ipnet.IP
					break
				}
			}
		}

		if ip == nil {
			panic("No IP Address found")
		} else {
			if err = s.service.ServerSet(ip.String()+s.port, true); err != nil {
				panic(err)
			} else {
				s.ip = ip.String()
			}
		}
	}
}

func (s *Server) StartServer() error {
	s.setServerInfo()

	channel := make(chan os.Signal, 1)
	signal.Notify(channel, syscall.SIGINT)

	go func() {
		<-channel
		if err := s.service.ServerSet(s.ip+s.port, false); err != nil {
			log.Println("Failed to set server info when server was closed", "err", err)
		}

		type ServerInfoEvent struct {
			IP     string
			Status bool
		}

		event := &ServerInfoEvent{
			IP:     s.ip + s.port,
			Status: false,
		}

		ch := make(chan kafka.Event)

		if value, err := json.Marshal(event); err != nil {
			log.Println("Failed to Marshal")
		} else if result, err := s.service.PublishEvent("chat", value, ch); err != nil {
			log.Println("Failed to send Event to kafka", "err", err)
		} else {
			// Send Event to kafka
			log.Println("Success to send Event to kafka", result)
		}

		os.Exit(1)
	}()

	log.Println("Go Server Starting...")
	return s.engine.Run(s.port)
}