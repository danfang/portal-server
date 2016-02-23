package model

import "github.com/jinzhu/gorm"

type Phone struct {
	gorm.Model
	User        User
	UserID      uint   `sql:"not null"`
	PhoneNumber string `sql:"not null"`
	Verified    bool   `sql:"not null; default false"`
}
