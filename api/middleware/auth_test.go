package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/model"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"portal-server/store"
)

const expectedResponse = "done"

var auth *gin.Engine
var authStore store.Store

func init() {
	gin.SetMode(gin.TestMode)
}

func createUser(uuid, token string, verified bool) {
	user := model.User{
		UUID:     uuid,
		Email:    uuid + "@portal.com",
		Verified: verified,
	}
	authStore.Users().CreateUser(&user)
	authStore.UserTokens().CreateToken(&model.UserToken{User: user, Token: token})
}

func init() {
	authStore = store.GetTestStore()

	createUser("1", "user_token_1", false)
	createUser("2", "user_token_2", true)

	auth = gin.New()
	gin.SetMode(gin.TestMode)
	auth.Use(AuthenticationMiddleware(authStore))
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
