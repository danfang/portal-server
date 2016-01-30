package main

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var api *gin.Engine

func init() {
	gin.SetMode(gin.TestMode)
	db, _ := gorm.Open("sqlite3", ":memory:")
	db.LogMode(false)
	api = API(&db)
}

func TestIndex(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1", nil)
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRegister(t *testing.T) {
	req, _ := http.NewRequest("POST", "/v1/register", bytes.NewBufferString(""))
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLogin(t *testing.T) {
	req, _ := http.NewRequest("POST", "/v1/login", bytes.NewBufferString(""))
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGoogleLogin(t *testing.T) {
	req, _ := http.NewRequest("POST", "/v1/login/google", bytes.NewBufferString(""))
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVerifyToken_MissingParameter(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/verify/", nil)
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestVerifyToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/verify/abc", nil)
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestAddDevice(t *testing.T) {
	req, _ := http.NewRequest("POST", "/v1/user/devices", nil)
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetDevices(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/user/devices", nil)
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestGetMessageHistory(t *testing.T) {
	req, _ := http.NewRequest("GET", "/v1/user/messages/history", nil)
	w := httptest.NewRecorder()
	api.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
