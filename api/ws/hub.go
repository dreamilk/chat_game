package ws

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"chat_game/log"
)

type Hub struct {
	m  map[string]*Client
	mu sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		m: make(map[string]*Client),
	}
}

var hub = NewHub()

func (h *Hub) Register(ctx context.Context, userID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.m[userID] = client

	log.Info(ctx, "register", zap.String("user_id", userID))

	// 确保在 client.Start() 之后再发送 welcome 消息
	go func() {
		client.send <- []byte("welcome")
	}()
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.m, userID)
}

func (h *Hub) Send(userID string, msg []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, ok := h.m[userID]
	if !ok {
		return errors.New("client not found")
	}

	client.send <- msg
	return nil
}
