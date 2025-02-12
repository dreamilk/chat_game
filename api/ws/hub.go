package ws

import (
	"context"
	"errors"
	"sync"

	"go.uber.org/zap"

	"chat_game/log"
)

type Hub struct {
	m       map[string]*Client
	mu      sync.Mutex
	receive chan HubMessage
}

type HubMessage struct {
	userID      string
	messageType int
	msg         []byte
}

func NewHub() *Hub {
	return &Hub{
		m:       make(map[string]*Client),
		receive: make(chan HubMessage, 1024),
	}
}

func init() {
	ctx := context.Background()
	go hub.Run(ctx)
}

var hub = NewHub()

func (h *Hub) Register(ctx context.Context, userID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.m[userID] = client

	log.Info(ctx, "register", zap.String("user_id", userID))

	client.send <- []byte("welcome")
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

func (h *Hub) Run(ctx context.Context) {
	for {
		msg := <-h.receive
		log.Info(ctx, "handle", zap.String("user_id", msg.userID), zap.String("msg", string(msg.msg)), zap.Int("message_type", msg.messageType))
	}
}
