package auth

import (
	"encoding/json"
	"github.com/danfang/portal-server/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const expectedResponse = "done"

var auth *gin.Engine
var authDb gorm.DB

func createUser(uuid, token string, verified bool) {
	user := model.User{
		UUID:     uuid,
		Email:    uuid + "@portal.com",
		Verified: verified,
	}
	authDb.Create(&user)
	authDb.Create(&model.UserToken{User: user, Token: token})
}

func init() {
	gin.SetMode(gin.TestMode)

	authDb, _ = gorm.Open("sqlite3", ":memory:")
	authDb.LogMode(false)
	authDb.CreateTable(&model.User{}, &model.UserToken{})

	createUser("1", "user_token_1", false)
	createUser("2", "user_token_2", true)

	auth = gin.New()
	auth.Use(AuthenticationMiddleware(&authDb))
	auth.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, expectedResponse)
	})
}

func TestAuthentication_MissingHeaders(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	expected, _ := json.Marshal(gin.H{"error": "missing_headers"})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestAuthentication_InvalidUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("X-USER-ID", "3")
	req.Header.Add("X-USER-TOKEN", "abcdef")
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	expected, _ := json.Marshal(gin.H{"error": "invalid_user_token"})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestAuthentication_InvalidUserToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("X-USER-ID", "1")
	req.Header.Add("X-USER-TOKEN", "incorrect_token")
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	expected, _ := json.Marshal(gin.H{"error": "invalid_user_token"})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestAuthentication_UnverifiedUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("X-USER-ID", "1")
	req.Header.Add("X-USER-TOKEN", "user_token_1")
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	expected, _ := json.Marshal(gin.H{"error": "account_not_verified"})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestAuthentication_VerifiedUser(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("X-USER-ID", "2")
	req.Header.Add("X-USER-TOKEN", "user_token_2")
	w := httptest.NewRecorder()
	auth.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, expectedResponse, w.Body.String())
}
