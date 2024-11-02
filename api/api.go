package api

import (
	"chat_game/handlers/room"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(eg *gin.Engine) {
	roomGroup := eg.Group("/room")
	roomGroup.GET("/list", room.List)
}
