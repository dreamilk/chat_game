package cmd

import (
	"context"

	"github.com/gin-gonic/gin"

	"chat_game/api"
	"chat_game/config"
	"chat_game/log"
)

func Server() {
	ctx := context.Background()

	log.Info(ctx, "server start")

	r := gin.New()
	r.Use(log.GinZap())

	api.RegisterRoute(r)
	go api.ServerRpc(ctx)

	r.Run(config.GetAppConfig().Port)
}
