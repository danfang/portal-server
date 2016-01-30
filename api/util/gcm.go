package util

import (
	"bytes"
	"encoding/json"
	"portal-server/api/errs"
	"io/ioutil"
	"net/http"
)

const (
	apiKey   = "AIzaSyAC4lW-Fb9tp12Un9LUiZNjw8ttVPQChPs"
	senderID = "1045304436932"
)

type notificationGroup struct {
	Operation string   `json:"operation"`
	KeyName   string   `json:"notification_key_name"`
	Key       string   `json:"notification_key,omitempty"`
	Tokens    []string `json:"registration_ids"`
}

type gcmResponse struct {
	Key   string `json:"notification_key"`
	Error string
}

// CreateNotificationGroup contacts Google GCM to create a new
// Cloud Messaging group, based on the given key and registration ID.
func CreateNotificationGroup(gcmClient *WebClient, keyName, registrationID string) (string, error) {
	data := notificationGroup{
		Operation: "create",
		KeyName:   keyName,
		Tokens:    []string{registrationID},
	}
	payload, err := json.Marshal(&data)
	if err != nil {
		return "", err
	}
	body, err := request(gcmClient, payload)
	if err != nil {
		return "", err
	}
	var res gcmResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}
	if res.Error != "" {
		return "", errs.GCMError(res.Error)
	}
	return res.Key, nil
}

// AddNotificationGroup contacts Google GCM to add a user device to an
// existing registration group.
func AddNotificationGroup(gcmClient *WebClient, keyName, key, registrationID string) error {
	data := notificationGroup{
		Operation: "add",
		KeyName:   keyName,
		Key:       key,
		Tokens:    []string{registrationID},
	}
	payload, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	body, err := request(gcmClient, payload)
	if err != nil {
		return err
	}
	var res gcmResponse
	if err := json.Unmarshal(body, &res); err != nil {
		return err
	}
	if res.Error != "" {
		return errs.GCMError(res.Error)
	}
	return nil
}

func request(gcmClient *WebClient, payload []byte) ([]byte, error) {
	req, err := http.NewRequest("POST", gcmClient.BaseURL, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+apiKey)
	req.Header.Set("project_id", senderID)
	res, err := gcmClient.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode >= 500 {
		return nil, errs.ErrGCMServiceUnavailable
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
