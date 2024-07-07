package network

import (
	"chat-server/service"
	"chat-server/types"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

type message struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Room    string `json:"room"`
	Time    int64  `json:"time"`
}

type client struct {
	Send   chan *message
	Room   *Room
	Name   string `json:"name"`
	Socket *websocket.Conn
}

func (c *client) Read() {
	// 클라이언트가 메시지를 읽는 함수
	defer c.Socket.Close()
	for {
		var msg *message
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			msg.Time = time.Now().Unix()
			msg.Name = c.Name

			c.Room.Forward <- msg
		}
	}
}

func (c *client) Write() {
	// 클라이언트가 메시지를 발행하는 함수
	defer c.Socket.Close()
	for msg := range c.Send {
		err := c.Socket.WriteJSON(msg)
		if err != nil {
			panic(err)
		}
	}
}

type Room struct {
	Forward chan *message // 수신되는 메시지를 보관하며 신규 메시지를 다른 클라이언트들에게 전달한다.
	Join    chan *client  // 소켓이 연결됐을 때 동작
	Leave   chan *client  // 소켓이 끊어졌을 때 동작
	Clients map[*client]bool
	service *service.Service
}

func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *client),
		Leave:   make(chan *client),
		Clients: make(map[*client]bool),
		service: service,
	}
}

func (r *Room) Run() {
	// Room 에 있는 모든 채널값들을 얻는 역할
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			close(client.Send)
			delete(r.Clients, client)
		case msg := <-r.Forward:

			go r.service.InsertChatting(msg.Name, msg.Message, msg.Room)

			for client := range r.Clients {
				client.Send <- msg
			}
		}
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  types.SocketBufferSize,
	WriteBufferSize: types.MessageBufferSize,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func (r *Room) ServeHttp(c *gin.Context) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("----serveHTTP:", err)
		return
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		log.Fatal("----auth Cookie failed:", err)
		return
	}

	client := &client{
		Socket: socket,
		Send:   make(chan *message, types.MessageBufferSize),
		Room:   r,
		Name:   userCookie.Value,
	}

	r.Join <- client
	defer func() { r.Leave <- client }()

	go client.Write()
	client.Read()
}
