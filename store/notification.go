package store

import "github.com/jinzhu/gorm"

type NotificationKeyStore interface {
}

type notificationKeyStore struct {
	*gorm.DB
}
