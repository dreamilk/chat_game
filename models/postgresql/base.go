package postgresql

import (
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	tgroupmessage "chat_game/models/postgresql/t_group_message"
	tmessage "chat_game/models/postgresql/t_message"
)

type DB struct {
	MessageDB      tmessage.MessageDB
	GroupMessageDB tgroupmessage.GroupMessageDB
}

func NewDB(dsn string) *DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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

	return &DB{
		MessageDB:      tmessage.NewMessageDB(db),
		GroupMessageDB: tgroupmessage.NewGroupMessageDB(db),
	}
}
