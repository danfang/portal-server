package access

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/api/errs"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/store"
	"testing"

	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var loginStore = store.GetTestStore()

func TestLoginEndpoing_InvalidEmail(t *testing.T) {
	input := map[string]string{
		"email":    "email",
		"password": "password",
	}
	w := testLogin(input)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidJSON.Error())
}

func TestLoginEndpoint_MissingFields(t *testing.T) {
	input := map[string]string{
		"email": "email@domain.com",
	}
	w := testLogin(input)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidJSON.Error())
}

func TestLoginEndpoint_NoSuchUser(t *testing.T) {
	input := map[string]string{
		"email":    "email@domain.com",
		"password": "my_password",
	}
	w := testLogin(input)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidLogin.Error())
}

func TestLoginEndpoint_BadPassowrd(t *testing.T) {
	createDefaultUser(loginStore, &passwordRegistration{
		Email:    "email@domain.com",
		Password: "my_password",
	})
	input := map[string]string{
		"email":    "email@domain.com",
		"password": "incorrect_password",
	}
	w := testLogin(input)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidLogin.Error())
}

func TestLoginEndpoint_Valid(t *testing.T) {
	createDefaultUser(loginStore, &passwordRegistration{
		Email:    "email2@domain.com",
		Password: "my_password",
	})
	input := map[string]string{
		"email":    "email2@domain.com",
		"password": "my_password",
	}
	w := testLogin(input)
	assert.Equal(t, 200, w.Code)
	assertValidLoginResponse(t, w)
}

func testLogin(input interface{}) *httptest.ResponseRecorder {
	// Create the router
	r := testutil.TestRouter(middleware.SetStore(loginStore))
	r.POST("/", LoginEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(body)))
	r.ServeHTTP(w, req)
	return w
}

func assertValidLoginResponse(t *testing.T, w *httptest.ResponseRecorder) {
	var res loginResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Regexp(t, "^[a-fA-F0-9]+$", res.UserToken)
	_, err := uuid.FromString(res.UserUUID)
	assert.NoError(t, err)
}
