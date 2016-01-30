package main

import (
	"encoding/json"
	"errors"
	"github.com/asaskevich/govalidator"
	"portal-server/model"
	"portal-server/model/types"
	"github.com/google/go-gcm"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"log"
)

// A GCMService handles upstream messages from a CloudConnectionService
// and sends appropriate responses downstream to clients. It also performs
// message validation and persistence.
type GCMService struct {
	Db  *gorm.DB
	CCS CloudConnectionServer
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
)

// MessagePayload is the message structure sent when a Portal client creates
// a new message and has broadcast it out to its device group.
type MessagePayload struct {
	MessageID string `json:"mid" valid:"required"`
	To        string `json:"to" valid:"required"`
	Status    string `json:"status" valid:"required,matches(started|sent|delivered|failed)"`
	Body      string `json:"body" valid:"required"`
	At        string `json:"at" valid:"required"`
}

// StatusPayload is the message structure sent when a Portal client updates
// the status of an existing message.
type StatusPayload struct {
	MessageID string `json:"mid" valid:"required"`
	Status    string `json:"status" valid:"required,matches(sent|delivered|failed)"`
	At        string `json:"at" valid:"required"`
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
	payload, ok := payload.(map[string]interface{})
	if !ok {
		return errors.New("invalid_payload_json")
	}
	bytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, result); err != nil {
		return err
	}
	if _, err := govalidator.ValidateStruct(result); err != nil {
		return err
	}
	return nil
}

func (s GCMService) recordMessage(cm gcm.CcsMessage, m MessagePayload) error {
	registrationID := cm.From
	var device model.Device
	if s.Db.Where(model.Device{
		RegistrationID: registrationID,
		State:          types.DeviceStateLinked.String(),
	}).First(&device).RecordNotFound() {
		return ErrUnregisteredDevice
	}
	message := &model.Message{
		UserID:    device.UserID,
		MessageID: m.MessageID,
		Status:    m.Status,
		To:        m.To,
		Body:      m.Body,
	}
	if err := s.Db.Create(message).Error; err != nil {
		return err
	}
	return nil
}
