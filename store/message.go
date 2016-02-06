package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type MessageStore interface {
	GetMessagesByUser(user *User, limit int) ([]Message, error)
	CreateMessage(proto *Message) error
}

type messageStore struct {
	*gorm.DB
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

func (db messageStore) CreateMessage(proto *Message) error {
	return db.Create(proto).Error
}
