package network

import (
	"chat-server/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type api struct {
	server *Server
}

func registerServer(server *Server) {
	api := &api{server: server}
	server.engine.GET("/room-list", api.roomList)
	server.engine.GET("/room", api.room)
	server.engine.GET("/enter-room", api.enterRoom)

	server.engine.POST("/make-room", api.makeRoom)

	r := NewRoom(server.service)
	go r.Run()
	server.engine.GET("/room-chat", r.ServeHttp)
}

func (api *api) roomList(ctx *gin.Context) {
	if res, err := api.server.service.RoomList(); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}

func (api *api) makeRoom(ctx *gin.Context) {
	var req types.BodyRoomReq

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if err := api.server.service.MakeRoom(req.Name); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, "Success")
	}
}

func (api *api) room(ctx *gin.Context) {
	var req types.FormRoomReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := api.server.service.Room(req.Name); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}

func (api *api) enterRoom(ctx *gin.Context) {
	var req types.FormRoomReq

	if err := ctx.ShouldBindQuery(&req); err != nil {
		response(ctx, http.StatusUnprocessableEntity, err.Error())
	} else if res, err := api.server.service.EnterRoom(req.Name); err != nil {
		response(ctx, http.StatusInternalServerError, err.Error())
	} else {
		response(ctx, http.StatusOK, res)
	}
}
