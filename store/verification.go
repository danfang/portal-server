package store

import (
	. "portal-server/model"

	"github.com/jinzhu/gorm"
)

type VerificationTokenStore interface {
	CreateToken(proto *VerificationToken) error
	FindToken(where *VerificationToken) (*VerificationToken, bool)
	FindDeletedToken(where *VerificationToken) (*VerificationToken, bool)
	DeleteToken(token *VerificationToken) error
	GetRelatedUser(token *VerificationToken) (*User, error)
	GetCount(where *VerificationToken) int
}

type verificationTokenStore struct {
	*gorm.DB
}

func (db verificationTokenStore) CreateToken(proto *VerificationToken) error {
	return db.Create(proto).Error
}

func (db verificationTokenStore) FindToken(where *VerificationToken) (*VerificationToken, bool) {
	var token VerificationToken
	if db.Where(where).First(&token).RecordNotFound() {
		return nil, false
	}
	return &token, true
}

func (db verificationTokenStore) FindDeletedToken(where *VerificationToken) (*VerificationToken, bool) {
	var token VerificationToken
	if db.Unscoped().Where(where).First(&token).RecordNotFound() {
		return nil, false
	}
	return &token, true
}

func (db verificationTokenStore) DeleteToken(token *VerificationToken) error {
	return db.Delete(token).Error
}

func (db verificationTokenStore) GetRelatedUser(token *VerificationToken) (*User, error) {
	var user User
	if err := db.Model(token).Related(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (db verificationTokenStore) GetCount(where *VerificationToken) int {
	var count int
	db.Model(&VerificationToken{}).Where(where).Count(&count)
	return count
}
