package main

import (
	"github.com/google/go-gcm"
	"log"
	"portal-server/store"
)

const (
	apiKey   = "AIzaSyAC4lW-Fb9tp12Un9LUiZNjw8ttVPQChPs"
	senderID = "1045304436932"
)

const (
	dbUser     = "portal_gcm"
	dbPassword = "password"
)

func init() {
	gcm.DebugMode = true
}

func main() {
	s := store.GetStore(dbUser, dbPassword)
	ccs := &GoogleCCS{senderID, apiKey}
	service := GCMService{Store: s, CCS: ccs}
	log.Fatal(service.CCS.Listen(service.OnMessageReceived, nil))
}
