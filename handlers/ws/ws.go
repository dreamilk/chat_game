package ws

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"chat_game/log"
	"chat_game/utils/opctx"
)

type MessageType string

const (
	MsgTypeRoomMsg MessageType = "to_room"
	MsgTypeUserMsg MessageType = "to_user"
)

type Message struct {
	Dst     string      `json:"dst"`
	MsgType MessageType `json:"msg_type"`
	Msg     string      `json:"msg"`
}

func HandleWs(ctx context.Context, h *Hub, userID string, msg []byte) {
	ctx = opctx.SetUserID(ctx, userID)

	var message Message
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Error(ctx, "unmarshal message", zap.Error(err), zap.String("msg", string(msg)), zap.String("user_id", userID))
		return
	}

	switch message.MsgType {
	case MsgTypeRoomMsg:
		if err := h.messageService.SendToRoom(ctx, message.Dst, []byte(message.Msg), h.Send); err != nil {
			log.Error(ctx, "send message to room", zap.Error(err), zap.String("room_id", message.Dst), zap.String("msg", message.Msg))
		}

	case MsgTypeUserMsg:
		if err := h.messageService.SendToUser(ctx, message.Dst, []byte(message.Msg), h.Send); err != nil {
			log.Error(ctx, "send message to user", zap.Error(err), zap.String("user_id", userID), zap.String("msg", message.Msg))
		}

	default:
		log.Error(ctx, "unknown message type", zap.String("msg_type", string(message.MsgType)))
	}
}
