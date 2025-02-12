package tgroupmessage

type GroupMessage struct {
	ID        int    `json:"id"`
	RoomID    int    `json:"room_id"`
	Sender    string `json:"sender"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
