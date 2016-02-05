package main

import (
	"github.com/google/go-gcm"
	"log"
	"portal-server/model"
)

const (
	apiKey   = "AIzaSyAC4lW-Fb9tp12Un9LUiZNjw8ttVPQChPs"
	senderID = "1045304436932"
)

const (
	dbUser     = "portal_gcm"
	dbName     = "portal"
	dbPassword = "password"
)

func init() {
	gcm.DebugMode = true
}

func main() {
	db := model.GetStore(dbUser, dbName, dbPassword)
	ccs := &GoogleCCS{senderID, apiKey}
	service := GCMService{db, ccs}
	log.Fatal(service.CCS.Listen(service.OnMessageReceived, nil))
}
