package store

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"portal-server/model"
)

func GetStore(dbUser, dbName, dbPassword string) Store {
	host := os.Getenv("POSTGRES_PORT_5432_TCP_ADDR")
	port := os.Getenv("POSTGRES_PORT_5432_TCP_PORT")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "5432"
	}

	connStr := fmt.Sprintf("user=%s dbname=%s host=%s port=%s password=%s sslmode=disable",
		dbUser, dbName, host, port, dbPassword)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return New(&db)
}

func GetTestStore() Store {
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.LogMode(false)
	db.CreateTable(&model.User{}, &model.VerificationToken{}, &model.LinkedAccount{}, &model.UserToken{},
		&model.NotificationKey{}, &model.Device{}, &model.Message{}, &model.Contact{}, &model.EncryptionKey{})
	return New(&db)
}
