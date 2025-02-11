package handlers

import (
	"github.com/gin-gonic/gin"
)

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Ping(ctx *gin.Context) {
	ctx.Writer.Write([]byte("pong"))
}
