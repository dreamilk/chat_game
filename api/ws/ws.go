package ws

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWs(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	ctx := context.Background()

	userID := c.Query("user_id")
	client := NewClient(hub, userID, conn)

	client.Start(ctx)
	hub.Register(ctx, userID, client)
}
