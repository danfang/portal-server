package store

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const dbName = "portal"

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

	if envPort != "" {
		port = envPort
	}
}

func GetStore(user, password string) Store {
	return New(GetDB(user, password))
}

func GetDB(user, password string) *gorm.DB {
	params := map[string]string{
		"dbname":   dbName,
		"host":     host,
		"port":     port,
		"user":     user,
		"password": password,
		"sslmode":  "disable",
	}
	var conn = ""
	for k, v := range params {
		conn += fmt.Sprintf("%s=%s ", k, v)
	}
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return &db
}
