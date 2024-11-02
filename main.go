package main

import (
	"chat_game/api"
	"chat_game/config"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.RegisterRoute(r)
	r.Run(config.GetAppConfig().Port)
}
