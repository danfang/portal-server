package store

import (
	"log"
	. "portal-server/model"

	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
)

func GetTestStore() Store {
	return New(GetTestDB())
}

func TeardownTestStore(s Store) {
	s.teardown()
}

func GetTestDB() *gorm.DB {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.LogMode(false)
	db.CreateTable(&User{}, &VerificationToken{}, &LinkedAccount{}, &UserToken{},
		&NotificationKey{}, &Device{}, &Message{}, &Contact{}, &ContactPhone{}, &EncryptionKey{})
	return &db
}

func TeardownTestDB(db *gorm.DB) {
	if _, valid := db.DB().Driver().(*sqlite3.SQLiteDriver); !valid {
		log.Fatalf("Teardown() should only be used in testing")
		return
	}
	db.DropTableIfExists(&User{}, &VerificationToken{}, &LinkedAccount{}, &UserToken{},
		&NotificationKey{}, &Device{}, &Message{}, &Contact{}, &ContactPhone{}, &EncryptionKey{})
}

func (s *store) teardown() {
	TeardownTestDB(s.db)
}
