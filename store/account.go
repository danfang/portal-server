package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type LinkedAccountStore interface {
	FindAccount(where *LinkedAccount) (*LinkedAccount, error)
	CreateAccount(proto *LinkedAccount) error
	GetRelatedUser(account *LinkedAccount) (*User, error)
	GetCount(where *LinkedAccount) int
}

type linkedAccountStore struct {
	*gorm.DB
}

func (db linkedAccountStore) FindAccount(where *LinkedAccount) (*LinkedAccount, error) {
	var account LinkedAccount
	if err := db.Where(where).First(&account).Error; err != nil {
		return nil, err
	}
	return &account, nil
}

func (db linkedAccountStore) CreateAccount(proto *LinkedAccount) error {
	return db.Create(proto).Error
}

func (db linkedAccountStore) GetRelatedUser(account *LinkedAccount) (*User, error) {
	var user User
	if err := db.Model(account).Related(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db linkedAccountStore) GetCount(where *LinkedAccount) int {
	var count int
	db.Where(where).Count(&count)
	return count
}
