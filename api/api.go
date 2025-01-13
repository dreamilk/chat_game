package api

import (
	"github.com/gin-gonic/gin"

	"chat_game/api/ws"
	"chat_game/handlers/room"
)

func RegisterRoute(eg *gin.Engine) {
	eg.GET("/ws", ws.ServeWs)

	roomGroup := eg.Group("/room")
	roomGroup.GET("/list", room.List)
}
