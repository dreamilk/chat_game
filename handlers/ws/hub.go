package ws

import (
	"context"
	"errors"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"chat_game/log"
	"chat_game/services/message"
	"chat_game/services/room"
)

type Hub struct {
	m              map[string]*Client
	mu             sync.Mutex
	receive        chan HubMessage
	roomService    room.RoomService
	messageService message.MessageService
}

type HubMessage struct {
	userID      string
	messageType int
	msg         []byte
}

var hub *Hub

func NewHub() *Hub {
	if hub != nil {
		return hub
	}

	hub = &Hub{
		m:              make(map[string]*Client),
		receive:        make(chan HubMessage, 1024),
		roomService:    room.NewRoomService(),
		messageService: message.NewMessageService(),
	}
	go hub.Run(context.Background())

	return hub
}

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

		if msg.messageType == websocket.TextMessage {
			nCtx := context.Background()

			HandleWs(nCtx, h, msg.userID, msg.msg)
		}
	}
}
