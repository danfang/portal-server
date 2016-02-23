package model

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	UUID      string `sql:"not null; type:uuid"`
	FirstName string
	LastName  string
	Email     string `sql:"not null; unique_index"`
	Password  string
	Verified  bool `sql:"not null; default false"`
}
