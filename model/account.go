package model

import "github.com/jinzhu/gorm"

const (
	LinkedAccountTypeGoogle = "google"
)

type LinkedAccount struct {
	gorm.Model
	User      User
	UserID    uint   `sql:"not null"`
	Type      string `sql:"not null"`
	AccountID string `sql:"not null"`
}
