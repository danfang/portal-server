package model

import "github.com/jinzhu/gorm"

type EncryptionKey struct {
	gorm.Model
	User   User
	UserID uint   `sql:"unique_index"`
	Key    string `sql:"not null"`
}
