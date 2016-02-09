package access

import (
	"portal-server/model"
	"portal-server/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

var userTokenStore = store.GetTestStore()

func TestCreateUserToken(t *testing.T) {
	user := &model.User{
		Email: "test@portal.com",
	}
	userTokenStore.Users().CreateUser(user)

	token, err := createUserToken(userTokenStore, user)
	assert.NoError(t, err)
	assert.Regexp(t, "^[a-fA-F0-9]+$", token)

	tokenFromDB, _ := userTokenStore.UserTokens().FindToken(&model.UserToken{Token: token})
	userFromDB, _ := userTokenStore.UserTokens().GetRelatedUser(tokenFromDB)
	assert.Equal(t, user.ID, userFromDB.ID)
	assert.Equal(t, user.Email, userFromDB.Email)
}
