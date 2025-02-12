package ws

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"chat_game/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	hub    *Hub
	userID string
	conn   *websocket.Conn
	send   chan []byte
}

func NewClient(hub *Hub, userID string, conn *websocket.Conn) *Client {
	return &Client{
		hub:    hub,
		userID: userID,
		conn:   conn,
		send:   make(chan []byte, 1024),
	}
}

func (c *Client) Start(ctx context.Context) {
	go c.read(ctx)
	go c.write(ctx)
}

func (c *Client) read(ctx context.Context) {
	defer func() {
		c.hub.Unregister(c.userID)
		c.conn.Close()
	}()

	for {
		messageType, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Info(ctx, "client error", zap.String("user_id", c.userID), zap.Error(err))
			}

			c.hub.receive <- HubMessage{
				userID:      c.userID,
				messageType: messageType,
				msg:         msg,
			}

			break
		}

		log.Info(ctx, "read", zap.String("msg", string(msg)))
		c.hub.receive <- HubMessage{
			userID:      c.userID,
			messageType: messageType,
			msg:         msg,
		}
	}
}

func (c *Client) write(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Info(ctx, "send channel closed")
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Error(ctx, "write error", zap.Error(err))
				return
			}
			log.Info(ctx, "write", zap.String("msg", string(msg)))

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error(ctx, "ping error", zap.Error(err))
				return
			}
		}
	}
}
