package tmessage

import "time"

type Message struct {
	ID        string        `json:"id" gorm:"default:uuid_generate_v4()"`
	Sender    string        `json:"sender"`
	Receiver  string        `json:"receiver"`
	Content   string        `json:"content"`
	Status    MessageStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type MessageStatus int

const (
	StatusUnread MessageStatus = 0
	StatusRead   MessageStatus = 1
)

const (
	ColumnID        = "id"
	ColumnSender    = "sender"
	ColumnReceiver  = "receiver"
	ColumnContent   = "content"
	ColumnStatus    = "status"
	ColumnCreatedAt = "created_at"
	ColumnUpdatedAt = "updated_at"
)

func (Message) TableName() string {
	return "t_message"
}
