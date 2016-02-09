package model

import "github.com/jinzhu/gorm"

const (
	DeviceStateLinked   = "linked"
	DeviceStateUnlinked = "unlinked"
)

const (
	DeviceTypePhone   = "phone"
	DeviceTypeChrome  = "chrome"
	DeviceTypeDesktop = "desktop"
)

type Device struct {
	gorm.Model
	User              User
	UserID            uint `sql:"not null"`
	NotificationKey   NotificationKey
	NotificationKeyID uint   `sql:"not null"`
	UUID              string `sql:"not null; type:uuid"`
	RegistrationID    string `sql:"not null; unique_index"`
	Name              string `sql:"not null"`
	Type              string `sql:"not null"`
	State             string `sql:"not null"`
}
