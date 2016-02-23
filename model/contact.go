package model

import "github.com/jinzhu/gorm"

type Contact struct {
	gorm.Model
	User         User
	UserID       uint           `sql:"not null"`
	UUID         string         `sql:"not null; type:uuid"`
	Name         string         `sql:"not null"`
	PhoneNumbers []ContactPhone `sql:"not null"`
}

type ContactPhone struct {
	gorm.Model
	ContactID uint   `sql:"not null"`
	Number    string `sql:"not null"`
	Name      string
}
