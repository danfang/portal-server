package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type MessageStore interface {
	FindMessage(where *Message) (*Message, bool)
	GetMessagesByUser(user *User, limit int) ([]Message, error)
	GetMessagesSince(user *User, messageID string) ([]Message, error)
	CreateMessage(proto *Message) error
	SaveMessage(message *Message) error
}

type messageStore struct {
	*gorm.DB
}

func (db messageStore) FindMessage(where *Message) (*Message, bool) {
	var message Message
	if db.Where(where).First(&message).RecordNotFound() {
		return nil, false
	}
	return &message, true
}

func (db messageStore) GetMessagesByUser(user *User, limit int) ([]Message, error) {
	var messages []Message
	if err := db.Where(&Message{
		UserID: user.ID,
	}).Order("id desc").Limit(limit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (db messageStore) GetMessagesSince(user *User, messageID string) ([]Message, error) {
	var message Message
	if err := db.Where(&Message{MessageID: messageID}).First(&message).Error; err != nil {
		return nil, err
	}
	var messages []Message
	if err := db.Where("id > ?", message.ID).Order("id desc").Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (db messageStore) CreateMessage(proto *Message) error {
	return db.Create(proto).Error
}

func (db messageStore) SaveMessage(message *Message) error {
	return db.Save(message).Error
}
