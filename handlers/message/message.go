package message

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chat_game/handlers"
	"chat_game/services/message"
)

var messageService message.MessageService

func init() {
	messageService = message.NewMessageService()
}

func List(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req struct {
		Receiver string `json:"receiver" form:"receiver" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	list, err := messageService.List(ctx, userID, req.Receiver)
	if err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "success",
		Data:    list,
	})
}

func ListRoom(ctx *gin.Context) {
	var req struct {
		RoomID string `json:"room_id"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	list, err := messageService.ListRoom(ctx, req.RoomID)
	if err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "success",
		Data:    list,
	})
}
