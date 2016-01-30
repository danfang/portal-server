package access

import (
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	"log"
)

func sendTokenToUser(email, token string) error {
	log.Println("Sending", token, "to", email)
	return nil
}

func hashPassword(password string, salt []byte) string {
	return hex.EncodeToString(pbkdf2.Key([]byte(password), salt, 4096, 48, sha256.New))
}
