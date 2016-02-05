package store

import "github.com/jinzhu/gorm"

type MessageStore interface {
}

type messageStore struct {
	*gorm.DB
}
