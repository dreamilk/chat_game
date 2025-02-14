package tmessage

import (
	"context"

	"gorm.io/gorm"
)

type MessageDB interface {
	List(ctx context.Context, userID string, friendID string) ([]Message, error)
	Insert(ctx context.Context, message Message) error
}

type MessageDBImpl struct {
	db *gorm.DB
}

func NewMessageDB(db *gorm.DB) MessageDB {
	if err := db.AutoMigrate(&Message{}); err != nil {
		panic(err)
	}

	return &MessageDBImpl{db: db}
}

func (m *MessageDBImpl) List(ctx context.Context, userID string, friendID string) ([]Message, error) {
	var list []Message
	err := m.db.WithContext(ctx).
		Where(ColumnSender+" = ? AND "+ColumnReceiver+" = ?", userID, friendID).
		Or(ColumnSender+" = ? AND "+ColumnReceiver+" = ?", friendID, userID).
		Order(ColumnCreatedAt + " DESC").
		Find(&list).Error
	return list, err
}

func (m *MessageDBImpl) Insert(ctx context.Context, message Message) error {
	return m.db.WithContext(ctx).Create(&message).Error
}
