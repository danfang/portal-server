package access

import (
	"github.com/danfang/portal-server/api/errs"
	"github.com/danfang/portal-server/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var verifyUserDB gorm.DB
var verifyUser model.User

func init() {
	verifyUserDB, _ = gorm.Open("sqlite3", ":memory:")
	verifyUserDB.CreateTable(&model.User{}, &model.VerificationToken{})
	verifyUserDB.LogMode(false)

	verifyUser = model.User{Email: "test@portal.com"}
	verifyUserDB.Create(&verifyUser)
}

func TestVerifyUserEndpoint_NoToken(t *testing.T) {
	w := testVerifyUser("")
	assert.Equal(t, 404, w.Code)
}

func TestVerifyUserEndpoint_BadToken(t *testing.T) {
	user := model.User{
		Email: "test_endpoint_bad_token@test.com",
	}
	verifyUserDB.Create(&user)
	token := model.VerificationToken{
		User:      user,
		Token:     "test_endpoint_bad_token",
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	verifyUserDB.Create(&token)
	w := testVerifyUser("invalid_token")
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidVerificationToken.Error())

	// Make sure token is not deleted
	var tokenCount int
	verifyUserDB.Model(model.VerificationToken{}).Where(model.VerificationToken{
		Token: "test_endpoint_bad_token",
	}).Count(&tokenCount)
	assert.Equal(t, 1, tokenCount)
}

func TestVerifyUserEndpoint_ExpiredToken(t *testing.T) {
	user := model.User{
		Email: "test_endpoint_expired_token@test.com",
	}
	verifyUserDB.Create(&user)
	token := model.VerificationToken{
		User:      user,
		Token:     "test_endpoint_expired_token",
		ExpiresAt: time.Now(),
	}
	verifyUserDB.Create(&token)
	w := testVerifyUser("test_endpoint_expired_token")
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrExpiredVerificationToken.Error())

	// Make sure token is not deleted
	var tokenCount int
	verifyUserDB.Model(model.VerificationToken{}).Where(model.VerificationToken{
		Token: "test_endpoint_expired_token",
	}).Count(&tokenCount)
	assert.Equal(t, 0, tokenCount)
}

func TestVerifyUserEndpoint_ValidToken(t *testing.T) {
	user := model.User{
		Email: "test_endpoint_valid_token@test.com",
	}
	verifyUserDB.Create(&user)
	token := model.VerificationToken{
		User:      user,
		Token:     "test_endpoint_valid_token",
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	verifyUserDB.Create(&token)
	w := testVerifyUser("test_endpoint_valid_token")
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"success":true}`, w.Body.String())

	// Make sure token is not deleted
	var tokenCount int
	verifyUserDB.Model(model.VerificationToken{}).Where(model.VerificationToken{
		Token: "test_endpoint_valid_token",
	}).Count(&tokenCount)
	assert.Equal(t, 0, tokenCount)
}

func TestCheckVerificationToken_NoSuchToken(t *testing.T) {
	_, err := checkVerificationToken(&verifyUserDB, "no_such_token")
	assert.EqualError(t, err, "invalid_verification_token")
}

func TestCheckVerificationToken_Expired(t *testing.T) {
	expiredToken := model.VerificationToken{
		User:      verifyUser,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
		Token:     "expired_token",
	}
	verifyUserDB.Create(&expiredToken)

	_, err := checkVerificationToken(&verifyUserDB, "expired_token")
	assert.EqualError(t, err, "expired_verification_token")

	var deletedToken model.VerificationToken
	verifyUserDB.Unscoped().Where("token = ?", "expired_token").First(&deletedToken)
	assert.NotNil(t, deletedToken.DeletedAt)
}

func TestCheckVerificationToken_NoUserToken(t *testing.T) {
	noUserToken := model.VerificationToken{
		UserID:    uint(404),
		ExpiresAt: time.Now().Add(time.Minute),
		Token:     "no_user_token",
	}
	verifyUserDB.Create(&noUserToken)

	_, err := checkVerificationToken(&verifyUserDB, "no_user_token")
	assert.EqualError(t, err, "invalid_verification_token")

	var deletedToken model.VerificationToken
	verifyUserDB.Unscoped().Where("token = ?", "no_user_token").First(&deletedToken)
	assert.NotNil(t, deletedToken.DeletedAt)
}

func TestCheckVerificationToken_ValidToken(t *testing.T) {
	token := model.VerificationToken{
		User:      verifyUser,
		ExpiresAt: time.Now().Add(time.Second),
		Token:     "token",
	}
	verifyUserDB.Create(&token)

	fromDB, err := checkVerificationToken(&verifyUserDB, "token")
	assert.NoError(t, err)
	assert.Equal(t, verifyUser.ID, fromDB.ID)

	var deletedToken model.VerificationToken
	verifyUserDB.Unscoped().Where("token = ?", "no_user_token").First(&deletedToken)
	assert.NotNil(t, deletedToken.DeletedAt)
}

func testVerifyUser(token string) *httptest.ResponseRecorder {
	// Create the router
	accessRouter := Router{&verifyUserDB, http.DefaultClient}
	r := gin.New()

	// Test the response
	r.GET("/:token", accessRouter.VerifyUserEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/"+token, nil)
	r.ServeHTTP(w, req)
	return w
}
