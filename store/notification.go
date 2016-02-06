package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type NotificationKeyStore interface {
	FindKey(where *NotificationKey) (*NotificationKey, bool)
	CreateKey(proto *NotificationKey) error
	GetRelatedUser(key *NotificationKey) (*User, error)
	GetCount(where *NotificationKey) int
}

type notificationKeyStore struct {
	*gorm.DB
}

func (db notificationKeyStore) FindKey(where *NotificationKey) (*NotificationKey, bool) {
	var key NotificationKey
	if db.Where(where).First(&key).RecordNotFound() {
		return nil, false
	}
	return &key, true
}

func (db notificationKeyStore) CreateKey(where *NotificationKey) error {
	return db.Create(where).Error
}

func (db notificationKeyStore) GetRelatedUser(key *NotificationKey) (*User, error) {
	var user User
	if err := db.Model(key).Related(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db notificationKeyStore) GetCount(where *NotificationKey) int {
	var count int
	db.Model(&NotificationKey{}).Where(where).Count(&count)
	return count
}
