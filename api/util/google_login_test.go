package util

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func expectLoginRequest(t *testing.T, idToken string) func(*http.Request) {
	return func(r *http.Request) {
		assert.Contains(t, r.RequestURI, "?id_token="+idToken)
	}
}

func TestCheckAUD(t *testing.T) {
	googleAUDs = []string{
		"valid_aud_1",
		"valid_aud_2",
		"valid_aud_3",
	}
	assert.True(t, checkAUD(&GoogleUser{Aud: "valid_aud_1"}))
	assert.True(t, checkAUD(&GoogleUser{Aud: "valid_aud_2"}))
	assert.True(t, checkAUD(&GoogleUser{Aud: "valid_aud_3"}))
	assert.False(t, checkAUD(&GoogleUser{Aud: "invalid_aud"}))
}

func TestGoogleLogin_InvalidJSONResponse(t *testing.T) {
	idToken := "my_id_token"
	requestTest := expectLoginRequest(t, idToken)
	server, client := TestHTTP(requestTest, 200, `{"bad_json"}`)
	defer server.Close()
	_, err := GetGoogleUser(client, idToken)
	assert.Error(t, err)
}

func TestGoogleLogin_Google400(t *testing.T) {
	idToken := "my_id_token"
	requestTest := expectLoginRequest(t, idToken)
	server, client := TestHTTP(requestTest, 400, `{}`)
	defer server.Close()
	_, err := GetGoogleUser(client, idToken)
	assert.EqualError(t, err, "invalid_google_id_token")
}

func TestGoogleLogin_GoogleNon200(t *testing.T) {
	idToken := "my_id_token"
	requestTest := expectLoginRequest(t, idToken)
	server, client := TestHTTP(requestTest, 500, `{}`)
	defer server.Close()
	_, err := GetGoogleUser(client, idToken)
	assert.EqualError(t, err, "google_oauth_unavailable")
}

func TestGoogleLogin_BadAUD(t *testing.T) {
	googleAUDs = []string{
		"valid_aud",
	}
	idToken := "my_id_token"
	requestTest := expectLoginRequest(t, idToken)
	mockResponse, _ := json.Marshal(map[string]string{
		"iss":            "https://accounts.google.com",
		"sub":            "110169484474386276334",
		"azp":            "invalid_aud",
		"aud":            "invalid_aud",
		"iat":            "1433978353",
		"exp":            "1433981953",
		"email":          "testuser@gmail.com",
		"email_verified": "true",
		"name":           "Test User",
		"picture":        "photo.jpg",
		"given_name":     "Test",
		"family_name":    "User",
		"locale":         "en",
	})
	server, client := TestHTTP(requestTest, 200, string(mockResponse))
	defer server.Close()
	_, err := GetGoogleUser(client, idToken)
	assert.EqualError(t, err, "invalid_google_id_token")
}

func TestGoogleLogin(t *testing.T) {
	googleAUDs = []string{
		"valid_aud",
	}
	mockResponse, _ := json.Marshal(map[string]string{
		"iss":            "https://accounts.google.com",
		"sub":            "110169484474386276334",
		"azp":            "valid_aud",
		"aud":            "valid_aud",
		"iat":            "1433978353",
		"exp":            "1433981953",
		"email":          "testuser@gmail.com",
		"email_verified": "true",
		"name":           "Test User",
		"picture":        "photo.jpg",
		"given_name":     "Test",
		"family_name":    "User",
		"locale":         "en",
	})
	idToken := "my_id_token"
	requestTest := expectLoginRequest(t, idToken)
	server, client := TestHTTP(requestTest, 200, string(mockResponse))
	defer server.Close()
	user, err := GetGoogleUser(client, idToken)
	assert.NoError(t, err)
	assert.Equal(t, user.Aud, "valid_aud")
	assert.Equal(t, user.GivenName, "Test")
	assert.Equal(t, user.FamilyName, "User")
	assert.Equal(t, user.Sub, "110169484474386276334")
	assert.Equal(t, user.EmailVerified, "true")
	assert.Equal(t, user.Email, "testuser@gmail.com")
	assert.Equal(t, user.Picture, "photo.jpg")
}
