package opctx

import (
	"context"

	"chat_game/utils/common"
)

func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, common.UserIDKey, userID)
}

func GetUserID(ctx context.Context) string {
	return ctx.Value(common.UserIDKey).(string)
}
