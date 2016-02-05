package model

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
)

var (
	host = "localhost"
	port = "5432"
)

func init() {
	envHost := os.Getenv("POSTGRES_PORT_5432_TCP_ADDR")
	envPort := os.Getenv("POSTGRES_PORT_5432_TCP_PORT")

	if envHost != "" {
		host = envHost
	}

	if port != "" {
		port = envPort
	}
}

func GetDB(dbUser, dbName, dbPassword string) *gorm.DB {
	connStr := fmt.Sprintf("user=%s dbname=%s host=%s port=%s password=%s sslmode=disable",
		dbUser, dbName, host, port, dbPassword)

	db, err := gorm.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return &db
}
