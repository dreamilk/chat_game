package cmd

import (
	"context"

	"github.com/gin-gonic/gin"

	"chat_game/api"
	"chat_game/config"
	"chat_game/log"
)

func Server() {
	log.Info(context.Background(), "server start")

	r := gin.New()
	r.Use(log.GinZap())

	api.RegisterRoute(r)
	r.Run(config.GetAppConfig().Port)
}
