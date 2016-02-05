package store

import "github.com/jinzhu/gorm"

type DeviceStore interface {
}

type deviceStore struct {
	*gorm.DB
}
