package wshub

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Hub struct {
	redis *redis.Client
}

func NewHub(redis *redis.Client) *Hub {
	return &Hub{redis: redis}
}

func (h *Hub) Register(ctx context.Context, userID string, wsAddr string) error {
	return h.redis.HSet(ctx, "users", userID, wsAddr).Err()
}

func (h *Hub) Unregister(ctx context.Context, userID string) error {
	return h.redis.HDel(ctx, "users", userID).Err()
}

func (h *Hub) Find(ctx context.Context, userID string) (string, error) {
	return h.redis.HGet(ctx, "users", userID).Result()
}
