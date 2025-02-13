package tuser

import (
	"context"

	"gorm.io/gorm"
)

type UserDB interface {
	Insert(ctx context.Context, user User) error
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id int) error
	One(ctx context.Context, id int) (User, error)
	List(ctx context.Context, offset, limit int) ([]User, error)
}

type UserDBImpl struct {
	db *gorm.DB
}

var _ UserDB = (*UserDBImpl)(nil)

func NewUserDB(db *gorm.DB) UserDB {
	return &UserDBImpl{db: db}
}

// One implements UserDB.
func (u *UserDBImpl) One(ctx context.Context, id int) (User, error) {
	var user User
	if err := u.db.First(&user, id).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

// List implements UserDB.
func (u *UserDBImpl) List(ctx context.Context, offset int, limit int) ([]User, error) {
	var users []User
	if err := u.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserDBImpl) Insert(ctx context.Context, user User) error {
	return u.db.Create(&user).Error
}

func (u *UserDBImpl) Update(ctx context.Context, user User) error {
	return u.db.Save(&user).Error
}

func (u *UserDBImpl) Delete(ctx context.Context, id int) error {
	return u.db.Delete(&User{ID: id}).Error
}
