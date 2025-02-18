package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"chat_game/log"
	"chat_game/utils/common"
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
	hub       *Hub
	userID    string
	conn      *websocket.Conn
	send      chan common.Msg
	stop      chan struct{}
	closeOnce sync.Once
}

func NewClient(hub *Hub, userID string, conn *websocket.Conn) *Client {
	return &Client{
		hub:    hub,
		userID: userID,
		conn:   conn,
		send:   make(chan common.Msg, 1024),
		stop:   make(chan struct{}),
	}
}

func (c *Client) Start(ctx context.Context) {
	go c.read(ctx)
	go c.write(ctx)
}

func (c *Client) Close(ctx context.Context) {
	c.closeOnce.Do(func() {
		close(c.stop)
		delete(c.hub.m, c.userID)
		c.conn.Close()
		close(c.send)
	})
}

func (c *Client) read(ctx context.Context) {
	defer func() {
		c.Close(ctx)
	}()

	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	c.conn.SetCloseHandler(func(code int, text string) error {
		log.Info(ctx, "connection closed",
			zap.String("user_id", c.userID),
			zap.Int("code", code),
			zap.String("text", text))
		return nil
	})

	for {
		select {
		case <-c.stop:
			return
		default:
			messageType, msg, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					log.Info(ctx, "client error", zap.String("user_id", c.userID), zap.Error(err))
				}
				return
			}

			log.Info(ctx, "read", zap.String("msg", string(msg)))
			c.hub.receive <- HubMessage{
				userID:      c.userID,
				messageType: messageType,
				msg:         msg,
			}
		}
	}
}

func (c *Client) write(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.Close(ctx)
	}()

	for {
		select {
		case <-c.stop:
			return
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Info(ctx, "send channel closed")
				return
			}

			msgBytes, err := json.Marshal(msg)
			if err != nil {
				log.Error(ctx, "marshal message", zap.Error(err))
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
				log.Error(ctx, "write error", zap.Error(err))
				return
			}
			log.Info(ctx, "write", zap.Any("msg", msg))

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Error(ctx, "ping error", zap.Error(err))
				return
			}
		}
	}
}
