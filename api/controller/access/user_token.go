package access

import (
	"crypto/rand"
	"encoding/hex"
	"portal-server/model"
	"time"

	"github.com/jinzhu/gorm"
)

func createUserToken(db *gorm.DB, user *model.User) (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	userToken := model.UserToken{
		User:      *user,
		ExpiresAt: time.Time{},
		Token:     hex.EncodeToString(token),
	}
	if err := db.Create(&userToken).Error; err != nil {
		return "", err
	}
	return userToken.Token, nil
}
