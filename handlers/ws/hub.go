package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"chat_game/log"
	"chat_game/rpc"
	client "chat_game/rpc/client"
	"chat_game/services/message"
	"chat_game/services/room"
	"chat_game/services/wshub"
	"chat_game/utils/common"
)

type Hub struct {
	m              map[string]*Client
	mu             sync.Mutex
	receive        chan HubMessage
	roomService    room.RoomService
	messageService message.MessageService
	rpcAddr        string
	msgHub         *wshub.Hub
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
		msgHub:         wshub.NewHub(),
	}
	go hub.Run(context.Background())

	return hub
}

func (h *Hub) Register(ctx context.Context, userID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.m[userID] = client

	log.Info(ctx, "register", zap.String("user_id", userID))

	client.send <- []byte("welcome " + userID)

	h.msgHub.Register(ctx, userID, h.rpcAddr)
}

func (h *Hub) Unregister(userID string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.m, userID)

	h.msgHub.Unregister(context.Background(), userID)
}

func (h *Hub) Send(userID string, msg []byte) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	c, ok := h.m[userID]
	if ok {
		c.send <- msg
		return nil
	}

	msgReq := common.Msg{}

	err := json.Unmarshal(msg, &msgReq)
	if err != nil {
		return err
	}

	ctx := context.Background()

	log.Info(ctx, "send message", zap.String("msg", string(msg)))

	rpcAddr, err := h.msgHub.Find(ctx, msgReq.Receiver)
	if err != nil {
		log.Error(ctx, "find msg rpc addr", zap.Error(err), zap.String("user_id", msgReq.Receiver))
		return err
	}

	msgRpcService, err := client.NewMsgServiceClient("tcp", rpcAddr)
	if err != nil {
		return err
	}

	_, err = msgRpcService.SendMessage(context.Background(), msgReq)
	if err != nil {
		return err
	}

	return nil
}

func (h *Hub) Run(ctx context.Context) {
	for {
		msg := <-h.receive
		log.Info(ctx, "handle", zap.String("user_id", msg.userID), zap.String("msg", string(msg.msg)), zap.Int("message_type", msg.messageType))

		if msg.messageType == websocket.TextMessage {
			nCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

			go func() {
				defer cancel()

				done := make(chan struct{})
				go func() {
					err := h.HandleWs(nCtx, msg.userID, msg.msg)
					if err != nil {
						log.Error(ctx, "handle ws", zap.String("user_id", msg.userID), zap.Error(err))
					}
					close(done)
				}()

				select {
				case <-done:
				case <-nCtx.Done():
					log.Error(ctx, "handle timeout",
						zap.String("user_id", msg.userID),
						zap.Error(nCtx.Err()))

					// for retry
					h.receive <- msg
				}
			}()
		}
	}
}

var _ rpc.MsgServer = (*Hub)(nil)

// SendMessage implements rpc.MsgServer.
func (h *Hub) SendMessage(req common.Msg, res *rpc.MsgResp) error {
	ctx := context.Background()

	msg, err := json.Marshal(req)
	if err != nil {
		return err
	}

	return h.HandleWs(ctx, req.Sender, msg)
}

func (h *Hub) SetRpcAddr(addr string) {
	h.rpcAddr = addr
}
