package tuser

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Age       int       `json:"age"`
	Sex       int       `json:"sex"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Sex int

const (
	SexMale   Sex = 1
	SexFemale Sex = 2
)

func (User) TableName() string {
	return "t_user"
}
