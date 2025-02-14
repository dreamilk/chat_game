package wshub

import (
	"context"

	"chat_game/config"
	"chat_game/models/redis"
)

type Hub struct {
	redisClient redis.Client
}

const (
	userRpcAddrKey = "user_rpc_addr"
)

func NewHub() *Hub {
	appConfig := config.GetAppConfig()
	redisClient := redis.NewRedis(appConfig.Redis.Addr, appConfig.Redis.User, appConfig.Redis.Password)

	return &Hub{redisClient: redisClient}
}

func (h *Hub) Register(ctx context.Context, userID string, wsAddr string) error {
	return h.redisClient.HSet(ctx, userRpcAddrKey, userID, wsAddr)
}

func (h *Hub) Unregister(ctx context.Context, userID string) error {
	return h.redisClient.HDel(ctx, userRpcAddrKey, userID)
}

func (h *Hub) Find(ctx context.Context, userID string) (string, error) {
	return h.redisClient.HGet(ctx, userRpcAddrKey, userID)
}
