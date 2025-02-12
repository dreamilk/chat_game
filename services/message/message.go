package message

import (
	"context"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"chat_game/config"
	tmessage "chat_game/models/mysql/t_message"
)

type MessageService interface {
	List(ctx context.Context, sender string, receiver string) ([]tmessage.Message, error)
	ListRoom(ctx context.Context, roomID string) ([]tmessage.Message, error)
	Insert(ctx context.Context, message tmessage.Message) error
}

type MessageServiceImpl struct {
	db *gorm.DB
}

func NewMessageService() MessageService {
	appConfig := config.GetAppConfig()

	db, err := gorm.Open(mysql.Open(appConfig.Mysql.Dsn), &gorm.Config{
		Logger: logger.New(log.New(log.Writer(), "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略 ErrRecordNotFound 错误
				Colorful:                  true,        // 禁用彩色打印
			},
		),
	})
	if err != nil {
		panic(err)
	}
	return &MessageServiceImpl{db: db}
}

func (m *MessageServiceImpl) List(ctx context.Context, sender string, receiver string) ([]tmessage.Message, error) {
	return tmessage.List(ctx, m.db, sender, receiver)
}

func (m *MessageServiceImpl) ListRoom(ctx context.Context, roomID string) ([]tmessage.Message, error) {
	return []tmessage.Message{}, nil
}

func (m *MessageServiceImpl) Insert(ctx context.Context, message tmessage.Message) error {
	return tmessage.Insert(ctx, m.db, message)
}
