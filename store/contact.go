package store

import (
	. "portal-server/model"

	"github.com/jinzhu/gorm"
)

type ContactStore interface {
	CreateContact(*Contact) error
	FindContact(where *Contact) (*Contact, bool)
	GetContactsByUser(user *User) ([]Contact, error)
}

type contactStore struct {
	*gorm.DB
}

func (db contactStore) CreateContact(proto *Contact) error {
	return db.Where(&Contact{
		UUID: proto.UUID,
	}).Assign(proto).FirstOrCreate(&Contact{}).Error
}

func (db contactStore) FindContact(where *Contact) (*Contact, bool) {
	var contact Contact
	if db.Where(where).Preload("PhoneNumbers").First(&contact).RecordNotFound() {
		return nil, false
	}
	return &contact, true
}

func (db contactStore) GetContactsByUser(user *User) ([]Contact, error) {
	var contacts []Contact
	if err := db.Where(&Contact{
		UserID: user.ID,
	}).Preload("PhoneNumbers").Find(&contacts).Error; err != nil {
		return nil, err
	}
	return contacts, nil
}
