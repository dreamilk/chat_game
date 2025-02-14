package message

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chat_game/handlers"
	tgroupmessage "chat_game/models/postgresql/t_group_message"
	tmessage "chat_game/models/postgresql/t_message"
	"chat_game/services/message"
)

var messageService message.MessageService

func init() {
	messageService = message.NewMessageService()
}

func List(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req struct {
		FriendID string `json:"friend_id" form:"friend_id" binding:"required,min=1"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	list, err := messageService.List(ctx, userID, req.FriendID)
	if err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	type resp struct {
		List  []tmessage.Message `json:"list"`
		Total int                `json:"total"`
	}
	r := resp{
		List:  list,
		Total: len(list),
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "success",
		Data:    r,
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

	type resp struct {
		List  []tgroupmessage.GroupMessage `json:"list"`
		Total int                          `json:"total"`
	}
	r := resp{
		List:  list,
		Total: len(list),
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "success",
		Data:    r,
	})
}
