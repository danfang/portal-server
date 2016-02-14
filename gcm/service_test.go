package main

import (
	"encoding/json"
	"errors"
	"portal-server/gcm/testutil"
	"portal-server/model"
	"portal-server/store"
	"testing"

	"github.com/franela/goblin"
	"github.com/google/go-gcm"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)

	g.Describe("GCM Service", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownStoreForTest(s)
		})

		g.It("Should correctly send a valid new message", func() {
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
			service := GCMService{s, ccs}
			messageID, err := service.sendMessage(&message)

			assert.NoError(t, err)
			assert.Equal(t, "message_id", messageID)
		})

		g.It("Should return a gcm_error on CCS failure", func() {
			message := gcm.XmppMessage{}
			ccs := testutil.TestCCS{
				XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
					assert.Equal(t, message, *m)
					return "", 400, errors.New("gcm_error")
				},
			}
			service := GCMService{s, ccs}
			messageID, err := service.sendMessage(&message)

			assert.EqualError(t, err, "gcm_error")
			assert.Equal(t, "", messageID)
		})

		g.It("Should be able to send an error message downstream", func() {
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
						"error":  ErrInvalidMessageType.Error(),
						"reason": errorReason,
					}, m.Data)
					return "message_id", 200, nil
				},
			}
			service := GCMService{s, ccs}
			service.errorMessage(registrationID, ErrInvalidMessageType, errorReason)
		})

		g.It("Should record a valid new message from downstream", func() {
			registrationID := "registration_id"
			messageID := "message_id"
			user := model.User{
				Email: "test@test.com",
			}
			s.Users().CreateUser(&user)
			s.Devices().CreateDevice(&model.Device{
				User:           user,
				RegistrationID: registrationID,
				Type:           model.DeviceTypePhone,
				State:          model.DeviceStateLinked,
			})
			ccs := testutil.TestCCS{
				XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
					t.Fail() // Should not have to send a failure message
					return "", 200, nil
				},
			}
			service := GCMService{s, ccs}
			payload, _ := json.Marshal(map[string]interface{}{
				"mid":    messageID,
				"status": "started",
				"at":     1351700038,
				"to":     "encrypted_phone_number",
				"body":   "encrypted_body",
			})
			service.OnMessageReceived(gcm.CcsMessage{
				From: registrationID,
				Data: map[string]interface{}{
					"type":    "message",
					"payload": string(payload),
				},
			})
			fromDB, _ := s.Messages().FindMessage(&model.Message{MessageID: messageID})
			assert.NotNil(t, fromDB)
			assert.Equal(t, "started", fromDB.Status)
			assert.Equal(t, "encrypted_phone_number", fromDB.To)
			assert.Equal(t, "encrypted_body", fromDB.Body)
		})

		g.It("Should not record a new message if the sending device is not found and send an error downstream", func() {
			registrationID := "unregistered_device"
			messageID := "unregistered_message_id"
			ccs := testutil.TestCCS{
				XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
					assert.Equal(t, m.Data["error"], "unregistered_device")
					return "", 200, nil
				},
			}
			service := GCMService{s, ccs}
			payload, _ := json.Marshal(map[string]interface{}{
				"mid":    messageID,
				"status": "started",
				"at":     1351700038,
				"to":     "encrypted_phone_number",
				"body":   "encrypted_body",
			})
			service.OnMessageReceived(gcm.CcsMessage{
				From: registrationID,
				Data: map[string]interface{}{
					"type":    "message",
					"payload": string(payload),
				},
			})
			_, found := s.Messages().FindMessage(&model.Message{MessageID: messageID})
			assert.False(t, found)
		})

		g.It("Should not record a new message and send an error on a message payload with a bad discriminator", func() {
			registrationID := "unregistered_device"
			messageID := "unregistered_message_id"
			ccs := testutil.TestCCS{
				XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
					assert.Equal(t, m.Data["error"], "invalid_message_type")
					return "", 200, nil
				},
			}
			service := GCMService{s, ccs}
			payload, _ := json.Marshal(map[string]interface{}{
				"mid":    messageID,
				"status": "started",
				"at":     1351700038,
				"to":     "encrypted_phone_number",
				"body":   "encrypted_body",
			})
			service.OnMessageReceived(gcm.CcsMessage{
				From: registrationID,
				Data: map[string]interface{}{
					"type":    "bad_discriminator",
					"payload": string(payload),
				},
			})
			_, found := s.Messages().FindMessage(&model.Message{MessageID: messageID})
			assert.False(t, found)
		})

		g.It("Should not record a new message and send an error on an invalid message payload", func() {
			registrationID := "unregistered_device"
			messageID := "unregistered_message_id"
			ccs := testutil.TestCCS{
				XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
					assert.Equal(t, m.Data["error"], "invalid_message_payload")
					return "", 200, nil
				},
			}
			service := GCMService{s, ccs}
			service.OnMessageReceived(gcm.CcsMessage{
				From: registrationID,
				Data: map[string]interface{}{
					"type":    "message",
					"payload": `{"bad_field": 0}`,
				},
			})
			_, found := s.Messages().FindMessage(&model.Message{MessageID: messageID})
			assert.False(t, found)
		})

		g.It("Should update a message on valid status message", func() {
			registrationID := "registration_id"
			messageID := "message_id"
			user := model.User{
				Email: "test@test.com",
			}
			s.Users().CreateUser(&user)
			s.Devices().CreateDevice(&model.Device{
				User:           user,
				RegistrationID: registrationID,
				Type:           model.DeviceTypePhone,
				State:          model.DeviceStateLinked,
			})
			s.Messages().CreateMessage(&model.Message{
				User:      user,
				MessageID: messageID,
				To:        "to",
				Body:      "body",
				Status:    "started",
			})
			ccs := testutil.TestCCS{
				XMPPFunc: func(m *gcm.XmppMessage) (string, int, error) {
					t.Fail() // Should not have to send a failure message
					return "", 200, nil
				},
			}
			service := GCMService{s, ccs}
			payload, _ := json.Marshal(map[string]interface{}{
				"mid":    messageID,
				"status": "sent",
				"at":     1351700038,
			})
			service.OnMessageReceived(gcm.CcsMessage{
				From: registrationID,
				Data: map[string]interface{}{
					"type":    "status",
					"payload": string(payload),
				},
			})
			fromDB, _ := s.Messages().FindMessage(&model.Message{MessageID: messageID})
			assert.NotNil(t, fromDB)
			assert.Equal(t, "sent", fromDB.Status)
		})
	})

	g.Describe("GCM Message payload marshalling", func() {
		g.It("Should marshall a new message json body into a MessagePayload struct", func() {
			payload, _ := json.Marshal(map[string]interface{}{
				"mid":    "message_id",
				"to":     "phone_number",
				"status": "started",
				"body":   "hello",
				"at":     1351700038,
			})
			var m MessagePayload
			err := getPayload(string(payload), &m)
			assert.NoError(t, err)
			assert.Equal(t, "message_id", m.MessageID)
			assert.Equal(t, "phone_number", m.To)
			assert.Equal(t, "started", m.Status)
			assert.Equal(t, "hello", m.Body)
			assert.Equal(t, 1351700038, m.At)
		})

		g.It("Should return an error for an invalid new message json body", func() {
			payload := "message_id"
			var m MessagePayload
			err := getPayload(payload, &m)
			assert.Error(t, err)
		})

		g.It("Should return an error when the json body fails field validation", func() {
			payload, _ := json.Marshal(map[string]interface{}{
				"mid": "message_id",
			})
			var m MessagePayload
			err := getPayload(string(payload), &m)
			assert.Error(t, err)
		})

		g.It("Should marshall a status json body into a StatusPayload struct", func() {
			payload, _ := json.Marshal(map[string]interface{}{
				"mid":    "message_id",
				"status": "sent",
				"at":     1351700038,
			})
			var m StatusPayload
			err := getPayload(string(payload), &m)
			assert.NoError(t, err)
			assert.Equal(t, "message_id", m.MessageID)
			assert.Equal(t, "sent", m.Status)
			assert.Equal(t, 1351700038, m.At)
		})

		g.It("Should return an error when the json body fails field validation", func() {
			payload := map[string]interface{}{
				"mid":    "message_id",
				"status": "bad_status",
				"at":     1351700038,
			}
			var m StatusPayload
			err := getPayload(payload, &m)
			assert.Error(t, err)
		})
	})
}
