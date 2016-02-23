package store

import (
	. "portal-server/model"

	"github.com/jinzhu/gorm"
)

type ContactStore interface {
	CreateContact(*Contact) error
	FindContact(where *Contact) (*Contact, bool)
}

type contactStore struct {
	*gorm.DB
}

func (db contactStore) CreateContact(proto *Contact) error {
	return db.Create(proto).Error
}

func (db contactStore) FindContact(where *Contact) (*Contact, bool) {
	var contact Contact
	if db.Where(where).First(&contact).RecordNotFound() {
		return nil, false
	}
	var phones []ContactPhone
	if err := db.Model(&contact).Related(&phones).Error; err != nil {
		return nil, false
	}
	contact.PhoneNumbers = phones
	return &contact, true
}
