package store

import (
	"github.com/jinzhu/gorm"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"
	"log"
	. "portal-server/model"
)

func GetTestStore() Store {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.LogMode(false)
	db.CreateTable(&User{}, &VerificationToken{}, &LinkedAccount{}, &UserToken{},
		&NotificationKey{}, &Device{}, &Message{}, &Contact{}, &EncryptionKey{})
	return New(&db)
}

func TeardownStoreForTest(s Store) {
	s.teardown()
}

func (s *store) teardown() {
	if _, valid := s.db.DB().Driver().(*sqlite3.SQLiteDriver); !valid {
		log.Fatalf("Teardown() should only be used in testing")
		return
	}
	s.db.DropTableIfExists(&User{}, &VerificationToken{}, &LinkedAccount{}, &UserToken{},
		&NotificationKey{}, &Device{}, &Message{}, &Contact{}, &EncryptionKey{})
}
