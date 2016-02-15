package main

import (
	"encoding/json"
	"errors"
	"log"
	"portal-server/model"

	"github.com/asaskevich/govalidator"
	"github.com/google/go-gcm"
	"github.com/satori/go.uuid"
	"portal-server/store"
)

// A GCMService handles upstream messages from a CloudConnectionService
// and sends appropriate responses downstream to clients. It also performs
// message validation and persistence.
type GCMService struct {
	Store store.Store
	CCS   CloudConnectionServer
}

// Message keys
const (
	discriminator = "type"
	payload       = "payload"
)

// Discriminator types
const (
	typeMessage = "message"
	typeStatus  = "status"
)

// Errors
var (
	ErrInvalidMessagePayload = errors.New("invalid_message_payload")
	ErrInvalidMessageType    = errors.New("invalid_message_type")
	ErrUnregisteredDevice    = errors.New("unregistered_device")
	ErrMessageNotFound       = errors.New("message_not_found")
)

// MessagePayload is the message structure sent when a Portal client creates
// a new message and has broadcast it out to its device group.
type MessagePayload struct {
	MessageID string `json:"mid" valid:"required,uuidv4"`
	To        string `json:"to" valid:"required"`
	Status    string `json:"status" valid:"required,matches(started|sent|delivered|failed)"`
	Body      string `json:"body" valid:"required"`
	At        int    `json:"at" valid:"required"`
}

// StatusPayload is the message structure sent when a Portal client updates
// the status of an existing message.
type StatusPayload struct {
	MessageID string `json:"mid" valid:"required,uuidv4"`
	Status    string `json:"status" valid:"required,matches(sent|delivered|failed)"`
	At        int    `json:"at" valid:"required"`
}

// OnMessageReceived handles all incoming GCM messages, performing
// validation and sending responses as necessary.
func (s GCMService) OnMessageReceived(cm gcm.CcsMessage) error {
	log.Printf("msg %v from %v\n", cm.Data, cm.From)
	d := cm.Data
	switch d[discriminator] {
	case typeMessage:
		var message MessagePayload
		if err := getPayload(d[payload], &message); err != nil {
			s.errorMessage(cm.From, ErrInvalidMessagePayload, err.Error())
			return nil
		}
		if err := s.recordMessage(cm, message); err == ErrUnregisteredDevice {
			s.errorMessage(cm.From, err, "device not found")
			return nil
		}
	case typeStatus:
		var message StatusPayload
		if err := getPayload(d[payload], &message); err != nil {
			s.errorMessage(cm.From, ErrInvalidMessagePayload, err.Error())
			return nil
		}
		err := s.updateMessage(cm, message)
		if err == ErrUnregisteredDevice {
			s.errorMessage(cm.From, err, "device not found")
			return nil
		}
		if err == ErrMessageNotFound {
			s.errorMessage(cm.From, err, "message not found")
			return nil
		}
	default:
		s.errorMessage(cm.From, ErrInvalidMessageType, "must be 'message' or 'status'")
	}
	return nil
}

func (s GCMService) sendMessage(m *gcm.XmppMessage) (string, error) {
	messageID, _, sendErr := s.CCS.SendXMPP(m)
	if sendErr != nil {
		return "", sendErr
	}
	return messageID, nil
}

func (s GCMService) errorMessage(to string, err error, reason string) {
	s.sendMessage(&gcm.XmppMessage{
		To:        to,
		MessageId: uuid.NewV4().String(),
		Data: map[string]interface{}{
			"error":  err.Error(),
			"reason": reason,
		},
	})
}

func getPayload(payload interface{}, result interface{}) error {
	byteString, ok := payload.(string)
	if !ok {
		return errors.New("invalid_payload_json")
	}
	if err := json.Unmarshal([]byte(byteString), result); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(result); err != nil {
		return err
	}
	return nil
}

func (s GCMService) recordMessage(cm gcm.CcsMessage, m MessagePayload) error {
	registrationID := cm.From
	device, found := s.Store.Devices().FindDevice(&model.Device{
		RegistrationID: registrationID,
		State:          model.DeviceStateLinked,
	})
	if !found {
		return ErrUnregisteredDevice
	}
	message := &model.Message{
		UserID:    device.UserID,
		MessageID: m.MessageID,
		Status:    m.Status,
		To:        m.To,
		Body:      m.Body,
	}
	return s.Store.Messages().CreateMessage(message)
}

func (s GCMService) updateMessage(cm gcm.CcsMessage, m StatusPayload) error {
	registrationID := cm.From
	device, found := s.Store.Devices().FindDevice(&model.Device{
		RegistrationID: registrationID,
		State:          model.DeviceStateLinked,
	})
	if !found {
		return ErrUnregisteredDevice
	}
	message, found := s.Store.Messages().FindMessage(&model.Message{UserID: device.UserID, MessageID: m.MessageID})
	if !found {
		return ErrMessageNotFound
	}
	message.Status = m.Status
	return s.Store.Messages().SaveMessage(message)
}
