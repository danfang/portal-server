package main

import (
	"fmt"
	"os"
	"portal-server/model"
)

const (
	dbUser     = "portal_db"
	dbName     = "portal"
	dbPassword = "password"
)

func validAction(action string) bool {
	switch action {
	case "drop", "create", "migrate":
		return true
	}
	return false
}

func main() {
	args := os.Args[1:]

	if len(args) != 1 || !validAction(args[0]) {
		fmt.Println("Usage:", os.Args[0], "[drop|create|migrate]")
		os.Exit(1)
	}

	fmt.Println("Connecting to db: ", dbName)
	db := model.GetStore(dbUser, dbName, dbPassword)
	db.LogMode(true)

	switch args[0] {
	case "drop":
		db.DropTable(
			&model.User{}, &model.VerificationToken{}, &model.LinkedAccount{}, &model.UserToken{},
			&model.NotificationKey{}, &model.Device{}, &model.Message{}, &model.Contact{}, &model.EncryptionKey{})

	case "create":
		db.CreateTable(
			&model.User{}, &model.VerificationToken{}, &model.LinkedAccount{}, &model.UserToken{},
			&model.NotificationKey{}, &model.Device{}, &model.Message{}, &model.Contact{}, &model.EncryptionKey{})
		db.Model(&model.LinkedAccount{}).AddUniqueIndex("idx_linked_account_type_account_id", "type", "account_id")

	case "migrate":
		db.AutoMigrate(
			&model.User{}, &model.VerificationToken{}, &model.LinkedAccount{}, &model.UserToken{},
			&model.NotificationKey{}, &model.Device{}, &model.Message{}, &model.Contact{}, &model.EncryptionKey{})
		db.Model(&model.LinkedAccount{}).AddUniqueIndex("idx_linked_account_type_account_id", "type", "account_id")
	}
}
