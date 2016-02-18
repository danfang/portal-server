package model

import "github.com/jinzhu/gorm"

type Contact struct {
	gorm.Model
	User    User
	UserID  uint   `sql:"not null"`
	UUID    string `sql:"not null; type:uuid"`
	Contact string `sql:"not null; type:text"`
}
