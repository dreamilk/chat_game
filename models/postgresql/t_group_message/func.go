package tgroupmessage

import (
	"context"

	"gorm.io/gorm"
)

type GroupMessageDB interface {
	Insert(ctx context.Context, message GroupMessage) error
}

type GroupMessageDBImpl struct {
	db *gorm.DB
}

func NewGroupMessageDB(db *gorm.DB) GroupMessageDB {
	if err := db.AutoMigrate(&GroupMessage{}); err != nil {
		panic(err)
	}

	return &GroupMessageDBImpl{db: db}
}

func (m *GroupMessageDBImpl) Insert(ctx context.Context, message GroupMessage) error {
	return m.db.WithContext(ctx).Create(&message).Error
}
