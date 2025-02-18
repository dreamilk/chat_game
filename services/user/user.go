package user

import (
	"context"

	"chat_game/config"
	"chat_game/models/redis"
)

type User struct {
	ID string `json:"id"`
}

const (
	userRpcAddrKey = "user_rpc_addr"
)

type UserService interface {
	List(ctx context.Context) ([]*User, error)
}

type userServiceImpl struct {
	redisClient redis.Client
}

var _ UserService = &userServiceImpl{}

func NewUserService() UserService {
	appConfig := config.GetAppConfig()
	redisClient := redis.NewRedis(appConfig.Redis.Addr, appConfig.Redis.User, appConfig.Redis.Password)

	return &userServiceImpl{
		redisClient: redisClient,
	}
}

func (s *userServiceImpl) List(ctx context.Context) ([]*User, error) {
	users, err := s.redisClient.HGetAll(ctx, userRpcAddrKey)
	if err != nil {
		return nil, err
	}

	userList := make([]*User, 0, len(users))
	for id := range users {
		userList = append(userList, &User{ID: id})
	}

	return userList, nil
}
