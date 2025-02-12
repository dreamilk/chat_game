package room

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"chat_game/config"
	"chat_game/handlers"
	"chat_game/models/redis"
	"chat_game/services/room"
)

var roomService room.RoomService

func init() {
	appConfig := config.GetAppConfig()

	redisClient := redis.NewRedis(appConfig.Redis.Addr, appConfig.Redis.User, appConfig.Redis.Password)

	roomService = room.NewRoomService(redisClient)
}

func List(ctx *gin.Context) {
	roomList, err := roomService.List(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Data:    roomList,
		Message: "ok",
	})
}

func Create(ctx *gin.Context) {
	ownerID := ctx.GetString("user_id")

	var req struct {
		RoomName string `json:"room_name" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	room, err := roomService.Create(ctx, ownerID, req.RoomName)
	if err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Data:    room,
		Message: "ok",
	})
}

func Join(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req struct {
		RoomID string `json:"room_id" form:"room_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	if err := roomService.Join(ctx, req.RoomID, userID); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "ok",
	})
}

func Leave(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req struct {
		RoomID string `json:"room_id" form:"room_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	if err := roomService.Leave(ctx, req.RoomID, userID); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "ok",
	})
}

func Delete(ctx *gin.Context) {
	userID := ctx.GetString("user_id")

	var req struct {
		RoomID string `json:"room_id" form:"room_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	if err := roomService.Delete(ctx, req.RoomID, userID); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Message: "ok",
	})
}

func Detail(ctx *gin.Context) {
	var req struct {
		RoomID string `json:"room_id" form:"room_id" binding:"required"`
	}

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	room, err := roomService.Detail(ctx, req.RoomID)
	if err != nil {
		ctx.JSON(http.StatusOK, handlers.Resp{
			Code:    -1,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, handlers.Resp{
		Code:    0,
		Data:    room,
		Message: "ok",
	})
}
