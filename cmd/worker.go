package cmd

import (
	"context"

	"go.uber.org/zap"

	"chat_game/log"
	"chat_game/services/receiver"
)

func Worker() {
	ctx := context.Background()

	addrs := []string{"localhost:9092"}

	receiver, err := receiver.NewReceiver(addrs)
	if err != nil {
		log.Error(ctx, "Failed to create receiver", zap.Error(err))
		return
	}

	receiver.Consume(ctx, "chat_game", func(ctx context.Context, message []byte) {
		log.Info(ctx, "Received message", zap.String("message", string(message)))
	})
}
