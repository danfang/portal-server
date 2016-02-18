package store

import (
	. "portal-server/model"

	"github.com/jinzhu/gorm"
)

type UserTokenStore interface {
	FindToken(where *UserToken) (*UserToken, bool)
	DeleteToken(token *UserToken) error
	CreateToken(token *UserToken) error
	GetRelatedUser(token *UserToken) (*User, error)
}

type userTokenStore struct {
	*gorm.DB
}

func (db userTokenStore) FindToken(where *UserToken) (*UserToken, bool) {
	var userToken UserToken
	if db.Where(where).First(&userToken).RecordNotFound() {
		return nil, false
	}
	return &userToken, true
}

func (db userTokenStore) DeleteToken(token *UserToken) error {
	return db.Delete(token).Error
}

func (db userTokenStore) CreateToken(token *UserToken) error {
	return db.Create(token).Error
}

func (db userTokenStore) GetRelatedUser(token *UserToken) (*User, error) {
	var user User
	if err := db.Model(token).Related(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
