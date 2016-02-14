package store

import (
	"fmt"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"time"
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
	}
	var conn = ""
	for k, v := range params {
		conn += fmt.Sprintf("%s=%s ", k, v)
	}
	conn += "sslmode=disable"

	// Connect to DB
	db, err := gorm.Open("postgres", conn)
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	// Ping DB
	if err := pingDatabase(&db); err != nil {
		log.Fatalf("Database ping attempts failed: %v\n", err)
	}
	return &db
}

func pingDatabase(db *gorm.DB) (err error) {
	for i := 0; i < 10; i++ {
		err = db.DB().Ping()
		if err == nil {
			return
		}
		time.Sleep(time.Second)
	}
	return
}
