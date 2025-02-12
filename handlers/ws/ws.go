package ws

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"chat_game/log"
)

type Message struct {
	User   string `json:"user"`
	RoomID string `json:"room_id"`
	Msg    string `json:"msg"`
}

func HandleWs(ctx context.Context, h *Hub, userID string, msg []byte) {
	var message Message
	if err := json.Unmarshal(msg, &message); err != nil {
		log.Error(ctx, "unmarshal message", zap.Error(err))
		return
	}
}
