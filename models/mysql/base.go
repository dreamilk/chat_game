package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	tuser "chat_game/models/mysql/t_user"
)

type DB struct {
	UserDB tuser.UserDB
}

func NewDB(mysqlDsn string) *DB {
	db, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return &DB{UserDB: tuser.NewUserDB(db)}
}
