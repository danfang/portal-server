package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type DeviceStore interface {
	DeviceCount(where *Device) int
	CreateDevice(proto *Device) error
}

type deviceStore struct {
	*gorm.DB
}

func (db deviceStore) DeviceCount(where *Device) int {
	var count int
	db.Where(where).Count(*count)
	return int
}

func (db deviceStore) CreateDevice(proto *Device) error {
	return db.Create(proto).Error
}
