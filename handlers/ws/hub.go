package ws

import (
	"context"
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
	rpcAddr, err := h.msgHub.Find(ctx, userID)
	if err == nil && rpcAddr != "" {
		log.Info(ctx, "already registered user", zap.String("user_id", userID), zap.String("rpc_addr", rpcAddr))

		h.msgHub.Unregister(ctx, userID, rpcAddr)
	}

	if c, ok := h.m[userID]; ok {
		c.Close(ctx)
	}

	h.m[userID] = client

	log.Info(ctx, "register", zap.String("user_id", userID))

	client.send <- common.Msg{
		MsgType:  common.MsgTypeSystem,
		Sender:   userID,
		Receiver: userID,
		Content:  "welcome " + userID,
	}

	h.msgHub.Register(ctx, userID, h.rpcAddr)
}

func (h *Hub) Send(receiver string, msg common.Msg) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	c, ok := h.m[receiver]
	if ok {
		c.send <- msg
		return nil
	}

	ctx := context.Background()

	log.Info(ctx, "send message", zap.Any("msg", msg))

	rpcAddr, err := h.msgHub.Find(ctx, receiver)
	if err != nil {
		log.Error(ctx, "find msg rpc addr", zap.Error(err), zap.String("receiver", receiver))
		return err
	}

	msgRpcService, err := client.NewMsgServiceClient("tcp", rpcAddr)
	if err != nil {
		return err
	}

	_, err = msgRpcService.SendMessage(context.Background(), common.WsHubMsg{
		Receiver: receiver,
		Msg:      msg,
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *Hub) Run(ctx context.Context) {
	go h.msgHub.Run(ctx, func(ctx context.Context, action *wshub.UserActionMsg) {
		if action.Src == h.rpcAddr {
			return
		}

		log.Info(ctx, "process user action", zap.Any("action", action))

		if action.Action == wshub.UserActionLogout {
			h.mu.Lock()
			defer h.mu.Unlock()

			if c, ok := h.m[action.UserID]; ok {
				c.Close(ctx)
			}
		} else if action.Action == wshub.UserActionLogin {
			h.mu.Lock()
			defer h.mu.Unlock()

			if c, ok := h.m[action.UserID]; ok {
				c.Close(ctx)
			}
		}
	})

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

func (h *Hub) Close(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, c := range h.m {
		h.msgHub.Unregister(ctx, c.userID, h.rpcAddr)
		c.Close(ctx)
	}
}

var _ rpc.MsgServer = (*Hub)(nil)

// SendMessage implements rpc.MsgServer.
func (h *Hub) SendMessage(req common.WsHubMsg, res *rpc.MsgResp) error {
	return h.Send(req.Receiver, req.Msg)
}

func (h *Hub) SetRpcAddr(addr string) {
	h.rpcAddr = addr
}
