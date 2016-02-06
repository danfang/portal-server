package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type DeviceStore interface {
	DeviceCount(where *Device) int
	CreateDevice(proto *Device) error
	GetAllLinkedDevices(user *User) ([]Device, error)
	GetRelatedUser(device *Device) (*User, error)
}

type deviceStore struct {
	*gorm.DB
}

func (db deviceStore) DeviceCount(where *Device) int {
	var count int
	db.Where(where).Count(&count)
	return count
}

func (db deviceStore) CreateDevice(proto *Device) error {
	return db.Create(proto).Error
}

func (db deviceStore) GetAllLinkedDevices(user *User) ([]Device, error) {
	var devices []Device
	if err := db.Where(Device{
		UserID: user.ID,
		State:  DeviceStateLinked,
	}).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (db deviceStore) GetRelatedUser(device *Device) (*User, error) {
	var user User
	if err := db.Model(device).Related(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
