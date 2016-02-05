package model

import "github.com/jinzhu/gorm"

type NotificationKey struct {
	gorm.Model
	User      User
	UserID    uint   `sql:"not null"`
	GroupName string `sql:"not null"`
	Key       string `sql:"not null"`
}
