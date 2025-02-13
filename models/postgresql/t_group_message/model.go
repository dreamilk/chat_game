package tgroupmessage

import "time"

type GroupMessage struct {
	ID        string    `json:"id" gorm:"default:uuid_generate_v4()"`
	RoomID    string    `json:"room_id"`
	Sender    string    `json:"sender"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (GroupMessage) TableName() string {
	return "t_group_message"
}
