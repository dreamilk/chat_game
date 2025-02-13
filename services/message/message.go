package message

import (
	"context"

	"chat_game/config"
	"chat_game/models/postgresql"
	tmessage "chat_game/models/postgresql/t_message"
)

type MessageService interface {
	List(ctx context.Context, sender string, receiver string) ([]tmessage.Message, error)
	ListRoom(ctx context.Context, roomID string) ([]tmessage.Message, error)
	Insert(ctx context.Context, message tmessage.Message) error
}

type MessageServiceImpl struct {
	db *postgresql.DB
}

func NewMessageService() MessageService {
	appConfig := config.GetAppConfig()
	db := postgresql.NewDB(appConfig.Postgres.Dsn)
	return &MessageServiceImpl{db: db}
}

func (m *MessageServiceImpl) List(ctx context.Context, sender string, receiver string) ([]tmessage.Message, error) {
	return m.db.MessageDB.List(ctx, sender, receiver)
}

func (m *MessageServiceImpl) ListRoom(ctx context.Context, roomID string) ([]tmessage.Message, error) {
	return []tmessage.Message{}, nil
}

func (m *MessageServiceImpl) Insert(ctx context.Context, message tmessage.Message) error {
	return m.db.MessageDB.Insert(ctx, message)
}
