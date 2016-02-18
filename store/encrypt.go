package store

import (
	. "portal-server/model"

	"github.com/jinzhu/gorm"
)

type EncryptionKeyStore interface {
	FindKey(where *EncryptionKey) (*EncryptionKey, bool)
	CreateKey(proto *EncryptionKey) error
	GetRelatedUser(key *EncryptionKey) (*User, error)
	GetCount(where *EncryptionKey) int
}

type encryptionKeyStore struct {
	*gorm.DB
}

func (db encryptionKeyStore) FindKey(where *EncryptionKey) (*EncryptionKey, bool) {
	var key EncryptionKey
	if db.Where(where).First(&key).RecordNotFound() {
		return nil, false
	}
	return &key, true
}

func (db encryptionKeyStore) CreateKey(proto *EncryptionKey) error {
	return db.Create(proto).Error
}

func (db encryptionKeyStore) GetRelatedUser(key *EncryptionKey) (*User, error) {
	var user User
	if err := db.Model(key).Related(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db encryptionKeyStore) GetCount(where *EncryptionKey) int {
	var count int
	db.Model(&EncryptionKey{}).Where(where).Count(&count)
	return count
}
