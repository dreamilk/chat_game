package wshub

import (
	"context"
	"encoding/json"

	"go.uber.org/zap"

	"chat_game/config"
	"chat_game/log"
	"chat_game/models/redis"
)

type Hub struct {
	redisClient redis.Client
}

const (
	userRpcAddrKey = "user_rpc_addr"
	userActionChan = "chan_user_action"
)

type UserAction int

const (
	UserActionLogin UserAction = iota
	UserActionLogout
)

type UserActionMsg struct {
	Action UserAction `json:"action"`
	UserID string     `json:"user_id"`
	Src    string     `json:"src"`
}

type HandleFunc func(ctx context.Context, action *UserActionMsg)

func NewHub() *Hub {
	appConfig := config.GetAppConfig()
	redisClient := redis.NewRedis(appConfig.Redis.Addr, appConfig.Redis.User, appConfig.Redis.Password)

	return &Hub{redisClient: redisClient}
}

func (h *Hub) Run(ctx context.Context, f HandleFunc) {
	if f == nil {
		return
	}

	pubSub := h.redisClient.Subscribe(ctx, userActionChan)
	defer pubSub.Close()

	ch := pubSub.Channel()
	for msg := range ch {
		userActionMsg := &UserActionMsg{}
		err := json.Unmarshal([]byte(msg.Payload), userActionMsg)
		if err != nil {
			log.Error(ctx, "unmarshal user action msg", zap.Error(err))
			continue
		}

		f(ctx, userActionMsg)
	}
}

func (h *Hub) Register(ctx context.Context, userID string, wsAddr string) error {
	userActionMsg := &UserActionMsg{
		Action: UserActionLogin,
		UserID: userID,
		Src:    wsAddr,
	}
	userActionMsgBytes, err := json.Marshal(userActionMsg)
	if err != nil {
		return err
	}

	err = h.redisClient.Publish(ctx, userActionChan, string(userActionMsgBytes))
	if err != nil {
		return err
	}

	return h.redisClient.HSet(ctx, userRpcAddrKey, userID, wsAddr)
}

func (h *Hub) Unregister(ctx context.Context, userID string, wsAddr string) error {
	userActionMsg := &UserActionMsg{
		Action: UserActionLogout,
		UserID: userID,
		Src:    wsAddr,
	}
	userActionMsgBytes, err := json.Marshal(userActionMsg)
	if err != nil {
		return err
	}

	err = h.redisClient.Publish(ctx, userActionChan, string(userActionMsgBytes))
	if err != nil {
		return err
	}

	return h.redisClient.HDel(ctx, userRpcAddrKey, userID)
}

func (h *Hub) Find(ctx context.Context, userID string) (string, error) {
	return h.redisClient.HGet(ctx, userRpcAddrKey, userID)
}

func (h *Hub) Publish(ctx context.Context, channel string, message string) error {
	return h.redisClient.Publish(ctx, channel, message)
}
