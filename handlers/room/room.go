package room

import (
	"chat_game/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Room struct {
	ID string `json:"id"`
}

func List(ctx *gin.Context) {
	var rooms []Room
	rooms = append(rooms, Room{
		ID: "abc1",
	})

	ctx.JSON(http.StatusOK, handlers.Return{
		Code: 0,
		Data: rooms,
	})
}
