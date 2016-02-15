package main

import (
	"log"
	"os"
	"portal-server/store"

	"github.com/google/go-gcm"
)

var (
	senderID = os.Getenv("GCM_SENDER_ID")
	apiKey   = os.Getenv("GCM_API_KEY")
)

var (
	dbName   = os.Getenv("DB_NAME")
	user     = os.Getenv("DB_GCM_USER")
	password = os.Getenv("DB_GCM_PASSWORD")
)

func init() {
	gcm.DebugMode = true
}

func main() {
	if senderID == "" || apiKey == "" {
		log.Fatalln("Missing GCM_SENDER_ID or GCM_API_KEY environment variables")
	}

	if dbName == "" || user == "" || password == "" {
		log.Fatalln("Missing DB_NAME, DB_GCM_USER, or DB_GCM_PASSWORD environment variables")
	}

	store := store.GetStore(dbName, user, password)
	ccs := &GoogleCCS{senderID, apiKey}
	service := GCMService{Store: store, CCS: ccs}
	log.Fatal(service.CCS.Listen(service.OnMessageReceived, nil))
}
