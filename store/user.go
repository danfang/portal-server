package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type UserStore interface {
	CreateUser(user *User) error
	SaveUser(user *User) error
	FindUser(where *User) (*User, bool)
	FindOrCreateUser(where *User, attrs *User) (*User, error)
	UserCount(where *User) int
	GetRelated(user *User, related interface{}) error
}

type userStore struct {
	*gorm.DB
}

func (db userStore) CreateUser(proto *User) error {
	return db.Create(proto).Error
}

func (db userStore) SaveUser(user *User) error {
	return db.Save(user).Error
}

func (db userStore) FindUser(where *User) (*User, bool) {
	var user User
	if db.Where(where).Find(&user).RecordNotFound() {
		return nil, false
	}
	return &user, true
}

func (db userStore) FindOrCreateUser(where *User, attrs *User) (*User, error) {
	var user User
	if err := db.Where(where).Attrs(attrs).FirstOrCreate(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db userStore) UserCount(where *User) int {
	var count int
	db.Where(where).Count(&count)
	return count
}

func (db userStore) GetRelated(user *User, related interface{}) error {
	return db.Model(user).Related(related).Error
}
