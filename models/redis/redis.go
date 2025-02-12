package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	c *redis.Client
}

type Client interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
	Del(ctx context.Context, key string) error

	HSet(ctx context.Context, key string, field string, value string) error
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, field string) error
	HExists(ctx context.Context, key string, field string) (bool, error)

	SAdd(ctx context.Context, key string, value string) error
	SMembers(ctx context.Context, key string) ([]string, error)
	SRem(ctx context.Context, key string, value string) error
	SIsMember(ctx context.Context, key string, value string) (bool, error)
}

var _ Client = &RedisClient{}

func NewRedis(addr string, user string, password string) *RedisClient {
	return &RedisClient{
		c: redis.NewClient(&redis.Options{
			Addr:     addr,
			Username: user,
			Password: password,
		}),
	}
}

func (r *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return r.c.Get(ctx, key).Result()
}

func (r *RedisClient) Set(ctx context.Context, key string, value string) error {
	return r.c.Set(ctx, key, value, 0).Err()
}

func (r *RedisClient) Del(ctx context.Context, key string) error {
	return r.c.Del(ctx, key).Err()
}

func (r *RedisClient) HSet(ctx context.Context, key string, field string, value string) error {
	return r.c.HSet(ctx, key, field, value).Err()
}

func (r *RedisClient) HGet(ctx context.Context, key string, field string) (string, error) {
	return r.c.HGet(ctx, key, field).Result()
}

func (r *RedisClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.c.HGetAll(ctx, key).Result()
}

func (r *RedisClient) HDel(ctx context.Context, key string, field string) error {
	return r.c.HDel(ctx, key, field).Err()
}

func (r *RedisClient) HExists(ctx context.Context, key string, field string) (bool, error) {
	return r.c.HExists(ctx, key, field).Result()
}

func (r *RedisClient) SAdd(ctx context.Context, key string, value string) error {
	return r.c.SAdd(ctx, key, value).Err()
}

func (r *RedisClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.c.SMembers(ctx, key).Result()
}

func (r *RedisClient) SRem(ctx context.Context, key string, value string) error {
	return r.c.SRem(ctx, key, value).Err()
}

func (r *RedisClient) SIsMember(ctx context.Context, key string, value string) (bool, error) {
	return r.c.SIsMember(ctx, key, value).Result()
}
