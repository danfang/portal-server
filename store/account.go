package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type LinkedAccountStore interface {
	FindAccount(*LinkedAccount) (*LinkedAccount, error)
	CreateAccount(*LinkedAccount) error
	GetRelatedUser(*LinkedAccount) (*User, error)
}

type linkedAccountStore struct {
	*gorm.DB
}

func (s linkedAccountStore) FindAccount(proto *LinkedAccount) (*LinkedAccount, error) {
	return nil, nil
}

func (s linkedAccountStore) CreateAccount(proto *LinkedAccount) error {
	return nil
}

func (s linkedAccountStore) GetRelatedUser(account *LinkedAccount) (*User, error) {
	return nil, nil
}
