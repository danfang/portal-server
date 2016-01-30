package access

import (
	"encoding/hex"
	"github.com/danfang/portal-server/model"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

var registerDB gorm.DB

func init() {
	registerDB, _ = gorm.Open("sqlite3", ":memory:")
	registerDB.LogMode(false)
	registerDB.CreateTable(&model.User{}, &model.VerificationToken{}, &model.UserToken{})
}

func TestCreateDefaultUser(t *testing.T) {
	body := passwordRegistration{
		Email:    "email@domain.com",
		Password: "password",
	}
	user, err := createDefaultUser(&registerDB, &body)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, "email@domain.com")
	assert.NotEmpty(t, user.Password)
	tokens := strings.Split(user.Password, ":")
	salt, _ := hex.DecodeString(tokens[1])
	assert.Equal(t, tokens[0], hashPassword("password", salt))

	// Duplicate user should fail
	_, err = createDefaultUser(&registerDB, &body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestCreateVerificationToken(t *testing.T) {
	user := model.User{
		Email:    "test@portal.com",
		Verified: false,
	}
	registerDB.Create(&user)

	token, err := createVerificationToken(&registerDB, &user)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Regexp(t, "^[a-fA-F0-9]+$", token)

	var tokenFromDB model.VerificationToken
	registerDB.Where("token = ?", token).First(&tokenFromDB)
	assert.Equal(t, token, tokenFromDB.Token)

	var userFromDB model.User
	registerDB.Model(&tokenFromDB).Related(&userFromDB)
	assert.Equal(t, user.ID, userFromDB.ID)
}
