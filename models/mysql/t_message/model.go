package tmessage

import "time"

type Message struct {
	ID        int       `json:"id"`
	Sender    string    `json:"sender"`
	Receiver  string    `json:"receiver"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

const (
	ColumnID        = "id"
	ColumnSender    = "sender"
	ColumnReceiver  = "receiver"
	ColumnContent   = "content"
	ColumnCreatedAt = "created_at"
	ColumnUpdatedAt = "updated_at"
)

func (Message) TableName() string {
	return "t_message"
}
