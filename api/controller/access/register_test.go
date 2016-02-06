package access

import (
	"encoding/hex"
	"portal-server/model"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"portal-server/store"
)

var registerStore = store.GetTestStore()

func TestCreateDefaultUser(t *testing.T) {
	body := passwordRegistration{
		Email:    "email@domain.com",
		Password: "password",
	}
	user, err := createDefaultUser(registerStore, &body)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, "email@domain.com")
	assert.NotEmpty(t, user.Password)
	tokens := strings.Split(user.Password, ":")
	salt, _ := hex.DecodeString(tokens[1])
	assert.Equal(t, tokens[0], hashPassword("password", salt))

	// Duplicate user should fail
	_, err = createDefaultUser(registerStore, &body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestCreateVerificationToken(t *testing.T) {
	user := model.User{
		Email:    "test@portal.com",
		Verified: false,
	}
	registerStore.Users().CreateUser(&user)

	token, err := createVerificationToken(registerStore, &user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Regexp(t, "^[a-fA-F0-9]+$", token)

	tokenFromDB, _ := registerStore.VerificationTokens().FindToken(&model.VerificationToken{Token: token})
	assert.Equal(t, token, tokenFromDB.Token)

	userFromDB, _ := registerStore.VerificationTokens().GetRelatedUser(tokenFromDB)
	assert.Equal(t, user.ID, userFromDB.ID)
}
