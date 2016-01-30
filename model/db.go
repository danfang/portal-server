package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func GetDB(dbUser, dbName, dbPassword string) *gorm.DB {
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

	return &db
}
