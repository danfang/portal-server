package model

import "github.com/jinzhu/gorm"

type Contact struct {
	gorm.Model   `json:"-"`
	User         User           `json:"-"`
	UserID       uint           `sql:"not null"                          json:"-"`
	UUID         string         `sql:"not null; type:uuid; unique_index" json:"cid"           valid:"required,uuid"`
	Name         string         `sql:"not null"                          json:"name"          valid:"required"`
	PhoneNumbers []ContactPhone `sql:"not null"                          json:"phone_numbers" valid:"required"`
}

type ContactPhone struct {
	gorm.Model `json:"-"`
	ContactID  uint   `sql:"not null" json:"-"`
	Number     string `sql:"not null" json:"number" valid:"required"`
	Type       string `sql:"not null" json:"type"   valid:"required"`
}
