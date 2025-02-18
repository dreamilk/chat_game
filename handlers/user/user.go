package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chat_game/handlers"
	"chat_game/services/user"
)

func List(ctx *gin.Context) {
	userService := user.NewUserService()
	users, err := userService.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "success",
		Data:    users,
	})
}
