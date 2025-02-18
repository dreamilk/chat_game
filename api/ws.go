package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"chat_game/handlers/ws"
	"chat_game/log"
)

var hub = ws.NewHub()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	ctx := context.Background()
	// 监听ctr c注销服务
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c

		hub.Close(ctx)

		log.Info(ctx, "server stop")
		os.Exit(0)
	}()
}

func ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	ctx := context.Background()

	userID := c.GetString("user_id")
	client := ws.NewClient(hub, userID, conn)

	client.Start(ctx)
	hub.Register(ctx, userID, client)
}
