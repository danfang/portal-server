package util

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"portal-server/api/errs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func expectRequest(t *testing.T, body interface{}) func(*http.Request) {
	return func(r *http.Request) {
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "key="+apiKey, r.Header.Get("Authorization"))
		assert.Equal(t, senderID, r.Header.Get("project_id"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		expected, _ := json.Marshal(body)
		actual, _ := ioutil.ReadAll(r.Body)
		assert.JSONEq(t, string(expected), string(actual))
	}
}

func TestGCM_CreateNotificationGroup(t *testing.T) {
	notificationKeyName := "notificationKeyName"
	registrationID := "registrationID"
	requestTest := expectRequest(t, map[string]interface{}{
		"operation":             "create",
		"notification_key_name": notificationKeyName,
		"registration_ids":      []string{registrationID},
	})
	mockResponse, _ := json.Marshal(map[string]string{
		"notification_key": "a_notification_key",
	})

	server, client := TestHTTP(requestTest, 200, string(mockResponse))
	defer server.Close()

	notificationKey, err := CreateNotificationGroup(client, notificationKeyName, registrationID)
	assert.NoError(t, err)
	assert.Equal(t, "a_notification_key", notificationKey)
}

func TestGCM_CreateNotificationGroup_GCMError(t *testing.T) {
	notificationKeyName := "notificationKeyName"
	registrationID := "registrationID"
	requestTest := expectRequest(t, map[string]interface{}{
		"operation":             "create",
		"notification_key_name": notificationKeyName,
		"registration_ids":      []string{registrationID},
	})
	mockResponse, _ := json.Marshal(map[string]string{
		"error": "google_is_down",
	})

	server, client := TestHTTP(requestTest, 200, string(mockResponse))
	defer server.Close()

	_, err := CreateNotificationGroup(client, notificationKeyName, registrationID)
	_, isGCMError := err.(errs.GCMError)
	assert.True(t, isGCMError)
	assert.EqualError(t, err, "google_is_down")
}

func TestGCM_CreateNotificationGroup_GoogleError(t *testing.T) {
	notificationKeyName := "notificationKeyName"
	registrationID := "registrationID"
	requestTest := expectRequest(t, map[string]interface{}{
		"operation":             "create",
		"notification_key_name": notificationKeyName,
		"registration_ids":      []string{registrationID},
	})

	server, client := TestHTTP(requestTest, 500, "")
	defer server.Close()

	_, err := CreateNotificationGroup(client, notificationKeyName, registrationID)
	_, isGCMError := err.(errs.GCMError)
	assert.True(t, isGCMError)
	assert.EqualError(t, err, "gcm_service_unavailable")
}

func TestGCM_AddNotificationGroup(t *testing.T) {
	notificationKeyName := "notificationKeyName"
	notificationKey := "notificationKey"
	registrationID := "registrationID"
	requestTest := expectRequest(t, map[string]interface{}{
		"operation":             "add",
		"notification_key_name": notificationKeyName,
		"notification_key":      notificationKey,
		"registration_ids":      []string{registrationID},
	})

	server, client := TestHTTP(requestTest, 200, "{}")
	defer server.Close()

	err := AddNotificationGroup(client, notificationKeyName, notificationKey, registrationID)
	assert.NoError(t, err)
}

func TestGCM_AddNotificationGroup_GCMError(t *testing.T) {
	notificationKeyName := "notificationKeyName"
	notificationKey := "notificationKey"
	registrationID := "registrationID"
	requestTest := expectRequest(t, map[string]interface{}{
		"operation":             "add",
		"notification_key_name": notificationKeyName,
		"notification_key":      notificationKey,
		"registration_ids":      []string{registrationID},
	})

	mockResponse, _ := json.Marshal(map[string]string{
		"error": "google_is_down",
	})
	server, client := TestHTTP(requestTest, 200, string(mockResponse))
	defer server.Close()

	err := AddNotificationGroup(client, notificationKeyName, notificationKey, registrationID)
	_, isGCMError := err.(errs.GCMError)
	assert.True(t, isGCMError)
	assert.EqualError(t, err, "google_is_down")
}

func TestGCM_AddNotificationGroup_GoogleError(t *testing.T) {
	notificationKeyName := "notificationKeyName"
	notificationKey := "notificationKey"
	registrationID := "registrationID"
	requestTest := expectRequest(t, map[string]interface{}{
		"operation":             "add",
		"notification_key_name": notificationKeyName,
		"notification_key":      notificationKey,
		"registration_ids":      []string{registrationID},
	})

	server, client := TestHTTP(requestTest, 500, "")
	defer server.Close()

	err := AddNotificationGroup(client, notificationKeyName, notificationKey, registrationID)
	_, isGCMError := err.(errs.GCMError)
	assert.True(t, isGCMError)
	assert.EqualError(t, err, "gcm_service_unavailable")
}
