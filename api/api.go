package api

import (
	"chat_game/handlers/room"

	"github.com/gin-gonic/gin"
)

func RegisterRoute(eg *gin.Engine) {
	g := eg.Group("/room")
	g.GET("/list", room.List)
}
