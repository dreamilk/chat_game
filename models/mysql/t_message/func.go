package tmessage

import (
	"context"

	"gorm.io/gorm"
)

func List(ctx context.Context, db *gorm.DB, sender string, receiver string) ([]Message, error) {
	var list []Message
	err := db.WithContext(ctx).Where(ColumnSender+" = ? AND "+ColumnReceiver+" = ?", sender, receiver).Find(&list).Error
	return list, err
}

func Insert(ctx context.Context, db *gorm.DB, message Message) error {
	return db.WithContext(ctx).Create(&message).Error
}
