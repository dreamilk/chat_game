package api

import (
	"github.com/gin-gonic/gin"

	"chat_game/auth"
	"chat_game/handlers"
	"chat_game/handlers/room"
)

func RegisterRoute(eg *gin.Engine) {
	eg.GET("/ws", ServeWs)
	eg.GET("/ping", handlers.Ping)

	roomGroup := eg.Group("/room")
	roomGroup.GET("/list", room.List)
	roomGroup.GET("/detail", room.Detail)
	roomGroup.POST("/create", auth.Auth(), room.Create)
	roomGroup.POST("/join", auth.Auth(), room.Join)
	roomGroup.POST("/leave", auth.Auth(), room.Leave)
	roomGroup.POST("/delete", auth.Auth(), room.Delete)
}
