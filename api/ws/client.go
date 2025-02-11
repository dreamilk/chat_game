package ws

import (
	"context"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"chat_game/log"
)

type Client struct {
	hub     *Hub
	userID  string
	conn    *websocket.Conn
	send    chan []byte
	receive chan []byte
}

func NewClient(hub *Hub, userID string, conn *websocket.Conn) *Client {
	return &Client{
		hub:     hub,
		userID:  userID,
		conn:    conn,
		send:    make(chan []byte, 1024),
		receive: make(chan []byte, 1024),
	}
}

func (c *Client) Start(ctx context.Context) {
	go c.read(ctx)
	go c.write(ctx)
}

func (c *Client) read(ctx context.Context) {
	defer func() {
		close(c.send)    // 先关闭 send channel
		c.conn.Close()   // 然后关闭连接
		close(c.receive) // 最后关闭 receive channel
		c.hub.Unregister(c.userID)
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(ctx, "read error", zap.Error(err))
			return
		}

		log.Info(ctx, "read", zap.String("msg", string(msg)))
		c.receive <- msg
	}
}

func (c *Client) write(ctx context.Context) {
	defer func() {
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				log.Info(ctx, "send channel closed")
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Error(context.Background(), "write error", zap.Error(err))
				return
			}
			log.Info(ctx, "write", zap.String("msg", string(msg)))
		case <-ctx.Done():
			return
		}
	}
}
