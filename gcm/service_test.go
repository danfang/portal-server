package main

import (
	"errors"
	"github.com/danfang/portal-server/gcm/testutil"
	"github.com/danfang/portal-server/model"
	"github.com/danfang/portal-server/model/types"
	"github.com/google/go-gcm"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDb gorm.DB

func init() {
	testDb, _ = gorm.Open("sqlite3", ":memory:")
	testDb.CreateTable(&model.User{}, &model.Device{}, &model.Message{})
}

func TestSendMessage(t *testing.T) {
	message := gcm.XmppMessage{
		To: "a friend",
		Data: map[string]interface{}{
			"key":   "key",
			"value": "value",
		},
	}
	ccs := testutil.TestCCS{
		XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
			assert.Equal(t, message, *m)
			return "message_id", 200, nil
		},
	}
	service := GCMService{&testDb, ccs}
	messageID, err := service.sendMessage(&message)

	assert.NoError(t, err)
	assert.Equal(t, "message_id", messageID)
}

func TestSendMessage_CcsError(t *testing.T) {
	message := gcm.XmppMessage{}
	ccs := testutil.TestCCS{
		XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
			assert.Equal(t, message, *m)
			return "", 400, errors.New("gcm_error")
		},
	}
	service := GCMService{&testDb, ccs}
	messageID, err := service.sendMessage(&message)

	assert.EqualError(t, err, "gcm_error")
	assert.Equal(t, "", messageID)
}

func TestErrorMessage(t *testing.T) {
	registrationID := "registration_id"
	errorReason := "a reason"
	ccs := testutil.TestCCS{
		XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
			// Check it was sent to the correct id
			assert.Equal(t, registrationID, m.To)

			// Check the message has a valid UUID
			_, err := uuid.FromString(m.MessageId)
			assert.NoError(t, err)

			// Check that the data has the errors
			assert.EqualValues(t, map[string]interface{}{
				"error":  "invalid_message_type",
				"reason": "a reason",
			}, m.Data)
			return "message_id", 200, nil
		},
	}
	service := GCMService{&testDb, ccs}
	service.errorMessage(registrationID, ErrInvalidMessageType, errorReason)
}

func TestGetPayload_MessagePayload(t *testing.T) {
	payload := map[string]interface{}{
		"mid":    "message_id",
		"to":     "phone_number",
		"status": "started",
		"body":   "hello",
		"at":     "1351700038",
	}
	var m MessagePayload
	err := getPayload(payload, &m)
	assert.NoError(t, err)
	assert.Equal(t, "message_id", m.MessageID)
	assert.Equal(t, "phone_number", m.To)
	assert.Equal(t, "started", m.Status)
	assert.Equal(t, "hello", m.Body)
	assert.Equal(t, "1351700038", m.At)
}

func TestGetPayload_MessagePayload_InvalidType(t *testing.T) {
	payload := "message_id"
	var m MessagePayload
	err := getPayload(payload, &m)
	assert.Error(t, err)
}

func TestGetPayload_MessagePayload_ValidationFailure(t *testing.T) {
	payload := map[string]interface{}{
		"mid": "message_id",
	}
	var m MessagePayload
	err := getPayload(payload, &m)
	assert.Error(t, err)
}

func TestGetPayload_StatusPayload(t *testing.T) {
	payload := map[string]interface{}{
		"mid":    "message_id",
		"status": "sent",
		"at":     "1351700038",
	}
	var m StatusPayload
	err := getPayload(payload, &m)
	assert.NoError(t, err)
	assert.Equal(t, "message_id", m.MessageID)
	assert.Equal(t, "sent", m.Status)
	assert.Equal(t, "1351700038", m.At)
}

func TestGetPayload_StatusPayload_ValidationFailure(t *testing.T) {
	payload := map[string]interface{}{
		"mid":    "message_id",
		"status": "bad_status",
		"at":     "1351700038",
	}
	var m StatusPayload
	err := getPayload(payload, &m)
	assert.Error(t, err)
}

func TestOnMessageReceived_ValidNewMessage(t *testing.T) {
	registrationID := "registration_id"
	messageID := "message_id"
	user := model.User{
		Email: "test@test.com",
	}
	testDb.Create(&user)
	testDb.Create(&model.Device{
		User:           user,
		RegistrationID: registrationID,
		Type:           types.DeviceTypePhone.String(),
		State:          types.DeviceStateLinked.String(),
	})
	ccs := testutil.TestCCS{
		XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
			t.Fail() // Should not have to send a failure message
			return "", 200, nil
		},
	}
	service := GCMService{&testDb, ccs}
	service.OnMessageReceived(gcm.CcsMessage{
		From: registrationID,
		Data: map[string]interface{}{
			"type": "message",
			"payload": map[string]interface{}{
				"mid":    messageID,
				"status": "started",
				"at":     "2015-06-09 08:00:00",
				"to":     "encrypted_phone_number",
				"body":   "encrypted_body",
			},
		},
	})
	var fromDB model.Message
	testDb.Where("message_id = ?", messageID).First(&fromDB)
	assert.NotNil(t, fromDB)
	assert.Equal(t, "started", fromDB.Status)
	assert.Equal(t, "encrypted_phone_number", fromDB.To)
	assert.Equal(t, "encrypted_body", fromDB.Body)
}

func TestOnMessageReceived_DeviceNotFound(t *testing.T) {
	registrationID := "unregistered_device"
	messageID := "unregistered_message_id"
	ccs := testutil.TestCCS{
		XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
			assert.Equal(t, m.Data["error"], "unregistered_device")
			return "", 200, nil
		},
	}
	service := GCMService{&testDb, ccs}
	service.OnMessageReceived(gcm.CcsMessage{
		From: registrationID,
		Data: map[string]interface{}{
			"type": "message",
			"payload": map[string]interface{}{
				"mid":    messageID,
				"status": "started",
				"at":     "2015-06-09 08:00:00",
				"to":     "encrypted_phone_number",
				"body":   "encrypted_body",
			},
		},
	})
	var count int
	testDb.Model(&model.Message{}).Where("message_id = ?", messageID).Count(&count)
	assert.Equal(t, 0, count)
}

func TestOnMessageReceived_BadDiscriminator(t *testing.T) {
	registrationID := "unregistered_device"
	messageID := "unregistered_message_id"
	ccs := testutil.TestCCS{
		XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
			assert.Equal(t, m.Data["error"], "invalid_message_type")
			return "", 200, nil
		},
	}
	service := GCMService{&testDb, ccs}
	service.OnMessageReceived(gcm.CcsMessage{
		From: registrationID,
		Data: map[string]interface{}{
			"type": "bad_discriminator",
			"payload": map[string]interface{}{
				"mid":    messageID,
				"status": "started",
				"at":     "2015-06-09 08:00:00",
				"to":     "encrypted_phone_number",
				"body":   "encrypted_body",
			},
		},
	})
	var count int
	testDb.Model(&model.Message{}).Where("message_id = ?", messageID).Count(&count)
	assert.Equal(t, 0, count)
}

func TestOnMessageReceived_BadPayload(t *testing.T) {
	registrationID := "unregistered_device"
	messageID := "unregistered_message_id"
	ccs := testutil.TestCCS{
		XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
			assert.Equal(t, m.Data["error"], "invalid_message_payload")
			return "", 200, nil
		},
	}
	service := GCMService{&testDb, ccs}
	service.OnMessageReceived(gcm.CcsMessage{
		From: registrationID,
		Data: map[string]interface{}{
			"type":    "message",
			"payload": map[string]interface{}{"bad_field": 0},
		},
	})
	var count int
	testDb.Model(&model.Message{}).Where("message_id = ?", messageID).Count(&count)
	assert.Equal(t, 0, count)
}
