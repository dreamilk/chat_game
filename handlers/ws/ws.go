package ws

import (
	"context"
	"encoding/json"
	"errors"

	"go.uber.org/zap"

	"chat_game/log"
	"chat_game/utils/common"
	"chat_game/utils/opctx"
)

func (h *Hub) HandleWs(ctx context.Context, userID string, msg []byte) error {
	ctx = opctx.SetUserID(ctx, userID)

	var message common.Msg
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Error(ctx, "unmarshal message", zap.Error(err), zap.String("msg", string(msg)), zap.String("user_id", userID))
		return err
	}

	switch message.MsgType {
	case common.MsgTypeRoom:
		if err := h.messageService.SendToRoom(ctx, message.RoomID, []byte(message.Content), h.Send); err != nil {
			log.Error(ctx, "send message to room", zap.Error(err), zap.String("room_id", message.RoomID), zap.String("msg", message.Content))
		}

	case common.MsgTypeUser:
		if err := h.messageService.SendToUser(ctx, message.Receiver, []byte(message.Content), h.Send); err != nil {
			log.Error(ctx, "send message to user", zap.Error(err), zap.String("user_id", message.Receiver), zap.String("msg", message.Content))
		}

	default:
		log.Error(ctx, "unknown message type", zap.String("msg_type", string(message.MsgType)))
		return errors.New("unknown message type")
	}

	return nil
}
