package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type UserToken struct {
	gorm.Model
	ExpiresAt time.Time
	User      User
	UserID    uint   `sql:"not null"`
	Token     string `sql:"not null"`
}
