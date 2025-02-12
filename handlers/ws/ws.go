package ws

import (
	"context"
	"encoding/json"
	"time"

	"go.uber.org/zap"

	"chat_game/log"
	tmessage "chat_game/models/mysql/t_message"
)

type MessageType string

const (
	MsgTypeRoomMsg MessageType = "to_room"
	MsgTypeUserMsg MessageType = "to_user"
)

type Message struct {
	User    string      `json:"user"`
	Dst     string      `json:"dst"`
	MsgType MessageType `json:"msg_type"`
	Msg     string      `json:"msg"`
}

func HandleWs(ctx context.Context, h *Hub, userID string, msg []byte) {
	var message Message
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Error(ctx, "unmarshal message", zap.Error(err), zap.String("msg", string(msg)), zap.String("user_id", userID))
		return
	}

	switch message.MsgType {
	case MsgTypeRoomMsg:
		room, err := h.roomService.Detail(ctx, message.Dst)
		if err != nil {
			log.Error(ctx, "get room detail", zap.Error(err), zap.String("room_id", message.Dst))
			return
		}

		for _, user := range room.Users {
			if err := h.Send(user, []byte(message.Msg)); err != nil {
				log.Error(ctx, "send message to user", zap.Error(err), zap.String("user_id", user), zap.String("msg", message.Msg))
			}
		}

	case MsgTypeUserMsg:
		if err := h.Send(message.Dst, []byte(message.Msg)); err != nil {
			log.Error(ctx, "send message to user", zap.Error(err), zap.String("user_id", message.Dst), zap.String("msg", message.Msg))
		}

		if err := h.messageService.Insert(ctx, tmessage.Message{
			Sender:    userID,
			Receiver:  message.Dst,
			Content:   message.Msg,
			CreatedAt: time.Now(),
		}); err != nil {
			log.Error(ctx, "insert message", zap.Error(err))
		}

	default:
		log.Error(ctx, "unknown message type", zap.String("msg_type", string(message.MsgType)))
	}
}
