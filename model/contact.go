package model

import "github.com/jinzhu/gorm"

type Contact struct {
	gorm.Model
	User     User
	UserID   uint   `sql:"not null"`
	Contacts string `sql:"type:jsonb"`
}
