package access

import (
	"portal-server/model"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"testing"
)

var userTokenDB gorm.DB

func init() {
	userTokenDB, _ = gorm.Open("sqlite3", ":memory:")
	userTokenDB.CreateTable(&model.User{}, &model.UserToken{})
}

func TestCreateUserToken(t *testing.T) {
	user := model.User{
		Email: "test@portal.com",
	}
	userTokenDB.Create(&user)

	token, err := createUserToken(&userTokenDB, &user)
	assert.NoError(t, err)
	assert.Regexp(t, "^[a-fA-F0-9]+$", token)

	var tokenFromDB model.UserToken
	var userFromDB model.User
	userTokenDB.Where(model.UserToken{Token: token}).First(&tokenFromDB)
	userTokenDB.Model(&tokenFromDB).Related(&userFromDB)
	assert.Equal(t, user.ID, userFromDB.ID)
	assert.Equal(t, user.Email, userFromDB.Email)
}
