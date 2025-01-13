package room

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chat_game/handlers"
)

type Room struct {
	ID string `json:"id"`
}

func List(ctx *gin.Context) {
	var rooms []Room
	rooms = append(rooms, Room{
		ID: "abc1",
	})

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code: 0,
		Data: rooms,
	})
}
