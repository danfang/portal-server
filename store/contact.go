package store

import "github.com/jinzhu/gorm"

type ContactStore interface {
}

type contactStore struct {
	*gorm.DB
}
