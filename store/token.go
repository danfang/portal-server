package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type UserTokenStore interface {
	FindToken(where *UserToken) (*UserToken, bool)
	DeleteToken(token *UserToken) error
	CreateToken(token *UserToken) error
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
