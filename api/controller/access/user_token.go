package access

import (
	"crypto/rand"
	"encoding/hex"
	"portal-server/model"
	"portal-server/store"
	"time"
)

func createUserToken(store store.Store, user *model.User) (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", err
	}
	userToken := &model.UserToken{
		User:      *user,
		ExpiresAt: time.Time{},
		Token:     hex.EncodeToString(token),
	}
	if err := store.UserTokens().CreateToken(userToken); err != nil {
		return "", err
	}
	return userToken.Token, nil
}
