package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
)

type User struct {
	gorm.Model
	UUID        string `sql:"not null; type:uuid;"`
	FirstName   string
	LastName    string
	Email       string `sql:"not null; unique_index"`
	Password    string
	Verified    bool `sql:"not null; default:false"`
	PhoneNumber string
}

type VerificationToken struct {
	gorm.Model
	ExpiresAt time.Time
	User      User
	UserID    uint   `sql:"not null"`
	Token     string `sql:"unique_index"`
}

type LinkedAccount struct {
	gorm.Model
	User      User
	UserID    uint   `sql:"not null"`
	Type      string `sql:"not null"`
	AccountID string `sql:"not null"`
}

type UserToken struct {
	gorm.Model
	ExpiresAt time.Time
	User      User
	UserID    uint   `sql:"not null"`
	Token     string `sql:"not null"`
}

type Device struct {
	gorm.Model
	User           User
	UserID         uint   `sql:"not null"`
	RegistrationID string `sql:"not null; unique_index"`
	Name           string `sql:"not null"`
	Type           string `sql:"not null"`
	State          string `sql:"not null"`
}

type NotificationKey struct {
	gorm.Model
	User      User
	UserID    uint   `sql:"not null"`
	GroupName string `sql:"not null"`
	Key       string `sql:"not null"`
}

type Message struct {
	gorm.Model
	User      User
	UserID    uint   `sql:"not null"`
	MessageID string `sql:"not null"`
	Status    string `sql:"not null"`
	To        string `sql:"not null"`
	Body      string `sql:"not null"`
}

type Contact struct {
	gorm.Model
	User     User
	UserID   uint   `sql:"not null"`
	Contacts string `sql:"type:jsonb"`
}

type EncryptionKey struct {
	gorm.Model
	User   User
	UserID uint   `sql:"unique_index"`
	Key    string `sql:"not null"`
}

// type DirectMessage struct {
//  gorm.Model
//  From        User
//  To          User
//  Message     string
