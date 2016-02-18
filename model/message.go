package model

import "github.com/jinzhu/gorm"

const (
	MessageStatusStarted   = "started"
	MessageStatusSent      = "sent"
	MessageStatusDelivered = "delivered"
	MessageStatusFailed    = "failed"
)

type Message struct {
	gorm.Model
	User      User
	UserID    uint   `sql:"not null"`
	MessageID string `sql:"not null"`
	Status    string `sql:"not null"`
	To        string `sql:"not null"`
	Body      string `sql:"type:text; not null"`
}
