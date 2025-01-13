package main

import (
	"github.com/gin-gonic/gin"

	"chat_game/api"
	"chat_game/config"
)

func main() {
	r := gin.Default()
	api.RegisterRoute(r)
	r.Run(config.GetAppConfig().Port)
}
