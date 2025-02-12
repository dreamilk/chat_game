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
		ctx.JSON(http.StatusInternalServerError, handlers.Resp{
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
	ownerID := ctx.Query("user_id")
	roomName := ctx.Query("room_name")

	room, err := roomService.Create(ctx, ownerID, roomName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, handlers.Resp{
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
}

func Leave(ctx *gin.Context) {
}

func Delete(ctx *gin.Context) {
}

func Detail(ctx *gin.Context) {
}
