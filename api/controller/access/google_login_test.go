package access

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/api/errs"
	"portal-server/api/util"
	"portal-server/model"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"portal-server/store"
)

var googleLoginStore = store.GetTestStore()

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGoogleLoginEndpoint_BadInput(t *testing.T) {
	input := map[string]string{
		"id_token": "",
	}
	w := testGoogleLogin(input, 200, "")
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidJSON.Error())
}

func TestGoogleLoginEndpoint_BadIDToken(t *testing.T) {
	input := map[string]string{
		"id_token": "token",
	}
	w := testGoogleLogin(input, 400, "{}")
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidGoogleIDToken.Error())
}

func TestGoogleLoginEndpoint_Google404(t *testing.T) {
	input := map[string]string{
		"id_token": "token",
	}
	w := testGoogleLogin(input, 404, "")
	assert.Equal(t, 500, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrGoogleOAuthUnavailable.Error())
}

func TestGoogleLoginEndpoint_GoogleEmailUnverified(t *testing.T) {
	input := map[string]string{
		"id_token": "token",
	}
	output := util.GoogleUser{
		Sub:           "1000",
		Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
		Email:         "test@google.com",
		EmailVerified: "false",
	}
	w := testGoogleLogin(input, 200, output)
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrGoogleAccountNotVerified.Error())
}

func TestGoogleLoginEndpoint_Valid(t *testing.T) {
	input := map[string]string{
		"id_token": "token",
	}
	output := util.GoogleUser{
		Sub:           "valid_user_sub",
		Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
		Email:         "test@google.com",
		EmailVerified: "true",
	}
	w := testGoogleLogin(input, 200, output)
	// Check login response
	assert.Equal(t, 200, w.Code)
	assertValidLoginResponse(t, w)

	// Check linked account is in DB
	linkedAccount, _ := googleLoginStore.LinkedAccounts().FindAccount(&model.LinkedAccount{
		AccountID: "valid_user_sub",
		Type:      model.LinkedAccountTypeGoogle,
	})

	assert.Equal(t, "valid_user_sub", linkedAccount.AccountID)

	// Check user is created
	user, _ := googleLoginStore.LinkedAccounts().GetRelatedUser(linkedAccount)
	assert.Equal(t, "test@google.com", user.Email)
	assert.True(t, user.Verified)
}

func TestGoogleLoginEndpoint_ExistingUser(t *testing.T) {
	user := model.User{
		UUID:     uuid.NewV4().String(),
		Email:    "test2@google.com",
		Verified: false,
		Password: "my_password_hash",
	}
	googleLoginStore.Users().CreateUser(&user)
	input := map[string]string{
		"id_token": "token",
	}
	output := util.GoogleUser{
		Sub:           "existing_user_sub",
		Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
		Email:         "test2@google.com",
		EmailVerified: "true",
	}
	w := testGoogleLogin(input, 200, output)
	// Check login response
	assert.Equal(t, 200, w.Code)
	assertValidLoginResponse(t, w)

	// Check linked account is in DB
	linkedAccount, _ := googleLoginStore.LinkedAccounts().FindAccount(&model.LinkedAccount{
		AccountID: "existing_user_sub",
		Type:      model.LinkedAccountTypeGoogle,
	})

	assert.Equal(t, "existing_user_sub", linkedAccount.AccountID)

	// Check user is created
	fromDB, _ := googleLoginStore.LinkedAccounts().GetRelatedUser(linkedAccount)
	assert.Equal(t, "test2@google.com", fromDB.Email)
	assert.True(t, fromDB.Verified)

	// Check that password login is disabled
	assert.Equal(t, "", fromDB.Password)
}

func TestGoogleLoginEndpoint_ExistingUserAndGoogleAccount(t *testing.T) {
	user := model.User{
		UUID:     uuid.NewV4().String(),
		Email:    "test3@google.com",
		Password: "some_password",
	}
	googleLoginStore.Users().CreateUser(&user)
	account := model.LinkedAccount{
		User:      user,
		AccountID: "existing_user_and_account_sub",
		Type:      model.LinkedAccountTypeGoogle,
	}
	googleLoginStore.LinkedAccounts().CreateAccount(&account)
	input := map[string]string{
		"id_token": "token",
	}
	output := util.GoogleUser{
		Sub:           "existing_user_and_account_sub",
		Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
		Email:         "test3@google.com",
		EmailVerified: "true",
	}
	w := testGoogleLogin(input, 200, output)
	// Check login response
	assert.Equal(t, 200, w.Code)
	assertValidLoginResponse(t, w)

	// Check linked account is in DB
	linkedAccountCount := googleLoginStore.LinkedAccounts().GetCount(&model.LinkedAccount{
		AccountID: "existing_user_and_account_sub",
		Type:      model.LinkedAccountTypeGoogle,
	})

	assert.Equal(t, 1, linkedAccountCount)
}

func TestCreateLinkedGoogleAccount(t *testing.T) {
	googleUser := util.GoogleUser{
		GivenName:     "Jon",
		FamilyName:    "Snow",
		EmailVerified: "true",
		Email:         "google@google.com",
		Sub:           "10000",
	}
	user, err := createLinkedGoogleAccount(googleLoginStore, &googleUser)
	assert.NoError(t, err)

	fromDB, _ := googleLoginStore.Users().FindUser(&model.User{Email: "google@google.com"})
	assert.Equal(t, fromDB.FirstName, "Jon")
	assert.Equal(t, fromDB.LastName, "Snow")
	assert.True(t, fromDB.Verified)

	var linkedAccount model.LinkedAccount
	googleLoginStore.Users().GetRelated(user, &linkedAccount)
	assert.Equal(t, linkedAccount.AccountID, "10000")
	assert.Equal(t, linkedAccount.Type, "google")
}

func TestCreateLinkedGoogleAccount_ExistingUser_NoLinkedAccount(t *testing.T) {
	original := model.User{
		Email:     "stannis@portal.com",
		FirstName: "Stannis",
		LastName:  "Baratheon",
		Verified:  false,
		Password:  "my_password",
	}

	googleLoginStore.Users().CreateUser(&original)

	googleUser := util.GoogleUser{
		GivenName:     "Stan",
		FamilyName:    "The Mannis",
		EmailVerified: "true",
		Email:         "stannis@portal.com",
		Sub:           "12345",
	}

	user, err := createLinkedGoogleAccount(googleLoginStore, &googleUser)
	assert.NoError(t, err)

	fromDB, _ := googleLoginStore.Users().FindUser(&model.User{Email: "stannis@portal.com"})

	assert.Equal(t, "Stannis", fromDB.FirstName)
	assert.Equal(t, "Baratheon", fromDB.LastName)

	// Check that the account is now verified and password disabled.
	assert.True(t, fromDB.Verified)
	assert.Equal(t, "", fromDB.Password)

	var linkedAccount model.LinkedAccount
	googleLoginStore.Users().GetRelated(user, &linkedAccount)
	assert.Equal(t, "12345", linkedAccount.AccountID)
	assert.Equal(t, model.LinkedAccountTypeGoogle, linkedAccount.Type)
}

func TestCreateLinkedGoogleAccount_ExistingUser_ExistingLinkedAccount(t *testing.T) {
	googleAccountID := "10101"

	original := model.User{
		Email: "existing@portal.com",
	}

	googleLoginStore.Users().CreateUser(&original)

	linkedAccount := model.LinkedAccount{
		User:      original,
		AccountID: googleAccountID,
		Type:      model.LinkedAccountTypeGoogle,
	}

	googleLoginStore.LinkedAccounts().CreateAccount(&linkedAccount)

	googleUser := util.GoogleUser{
		Sub:   googleAccountID,
		Email: "otherEmail@otherDomain.com",
	}

	// Make sure no data is modified
	user, err := createLinkedGoogleAccount(googleLoginStore, &googleUser)
	assert.NoError(t, err)
	assert.Equal(t, original.ID, user.ID)
	assert.Equal(t, original.Email, user.Email)

	// Make sure no new linked account is created.
	count := googleLoginStore.LinkedAccounts().GetCount(&model.LinkedAccount{
		AccountID: googleAccountID,
		Type:      model.LinkedAccountTypeGoogle,
	})
	assert.Equal(t, 1, count)
}

func testGoogleLogin(input interface{}, googleResponseCode int, googleResponseBody interface{}) *httptest.ResponseRecorder {
	// Setup mock Google server/client
	output, _ := json.Marshal(googleResponseBody)
	server, client := util.TestHTTP(func(*http.Request) {}, googleResponseCode, string(output))
	defer server.Close()

	// Modify the OAuth endpoint
	googleOAuthEndpoint = client.BaseURL

	// Create the router based on the db and Mock client
	accessRouter := Router{googleLoginStore, client.HTTPClient}
	r := gin.New()

	// Test the response
	r.POST("/", accessRouter.GoogleLoginEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(body)))
	r.ServeHTTP(w, req)
	return w
}
