package store

import (
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	. "portal-server/model"
)

func GetTestStore() Store {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.LogMode(false)
	db.CreateTable(&User{}, &VerificationToken{}, &LinkedAccount{}, &UserToken{},
		&NotificationKey{}, &Device{}, &Message{}, &Contact{}, &EncryptionKey{})
	return New(&db)
}
