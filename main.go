package main

import (
	"chat_game/api"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	api.RegisterRoute(r)
	r.Run(":8001")
}
