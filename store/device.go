package store

import (
	"github.com/jinzhu/gorm"
	. "portal-server/model"
)

type DeviceStore interface {
	FindDevice(where *Device) (*Device, bool)
	DeleteDevice(device *Device) error
	DeviceCount(where *Device) int
	CreateDevice(proto *Device) error
	GetAllLinkedDevices(user *User) ([]Device, error)
	GetRelatedUser(device *Device) (*User, error)
	GetRelatedKey(device *Device) (*NotificationKey, error)
}

type deviceStore struct {
	*gorm.DB
}

func (db deviceStore) FindDevice(where *Device) (*Device, bool) {
	var device Device
	if db.Where(where).First(&device).RecordNotFound() {
		return nil, false
	}
	return &device, true
}

func (db deviceStore) DeleteDevice(device *Device) error {
	return db.Delete(device).Error
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

func (db deviceStore) GetRelatedKey(device *Device) (*NotificationKey, error) {
	var key NotificationKey
	if err := db.Model(device).Related(&key).Error; err != nil {
		return nil, err
	}
	return &key, nil
}
