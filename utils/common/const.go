package common

const (
	UserIDKey = "user_id"
)

type MsgType string

const (
	MsgTypeRoom MsgType = "room"
	MsgTypeUser MsgType = "user"
)

type Msg struct {
	MsgType  MsgType `json:"msg_type"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	RoomID   string  `json:"room_id"`
	Content  string  `json:"content"`
}
