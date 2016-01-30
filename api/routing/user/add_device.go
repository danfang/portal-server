package user

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/danfang/portal-server/api/errs"
	"github.com/danfang/portal-server/api/routing"
	"github.com/danfang/portal-server/api/util"
	"github.com/danfang/portal-server/model"
	"github.com/danfang/portal-server/model/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

const gcmEndpoint = "https://android.googleapis.com/gcm/notification"

// AddDevice is a JSON structure for registering a GCM device.
//
// swagger:parameters addDevice
type AddDevice struct {
	// in: body
	// required: true
	Body addDevice `json:"add_device"`
}

type addDevice struct {
	// required: true
	RegistrationID string `json:"registration_id" valid:"required"`

	// required: true
	Name string `json:"name" valid:"required"`

	// required: true
	// pattern: (phone,chrome,desktop)
	Type types.DeviceType `json:"type" valid:"required,matches(phone,chrome,desktop)"`
}

// AddDeviceResponse contains the encryption and
// notifica†ion keys for a new GCM client.
//
// swagger:response addDevice
type AddDeviceResponse struct {
	// in: body
	Body addDeviceResponse `json:"add_device"`
}

type addDeviceResponse struct {
	EncryptionKey   string `json:"encryption_key"`
	NotificationKey string `json:"notification_key"`
}

// AddDeviceEndpoint handles a POST request to register or add a new
// messaging device.
func (r Router) AddDeviceEndpoint(c *gin.Context) {
	userID := c.MustGet("userID").(uint)

	var body addDevice
	if !routing.ValidateJSON(c, &body) {
		return
	}

	tx := r.Db.Begin()

	var count int
	if tx.Model(model.Device{}).Where(model.Device{
		RegistrationID: body.RegistrationID,
	}).Count(&count); count >= 1 {
		c.JSON(http.StatusBadRequest, routing.RenderError(errs.ErrDuplicateDeviceToken))
		return
	}

	device, err := createDevice(tx, userID, &body)
	if err != nil {
		tx.Rollback()
		routing.InternalServiceError(c, err)
		return
	}

	gcmClient := &util.WebClient{gcmEndpoint, r.HTTPClient}

	notificationKey, err := createNotificationKey(tx, gcmClient, userID, device.RegistrationID)
	if err, isGCMError := err.(errs.GCMError); isGCMError {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, routing.DetailError{
			Error:  errs.ErrUnableToRegisterDevice.Error(),
			Reason: err.Error(),
		})
		return
	}
	if err != nil {
		tx.Rollback()
		routing.InternalServiceError(c, err)
		return
	}

	encryptionKey, err := createEncryptionKey(tx, userID)
	if err != nil {
		tx.Rollback()
		routing.InternalServiceError(c, err)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, addDeviceResponse{
		EncryptionKey:   encryptionKey.Key,
		NotificationKey: notificationKey.Key,
	})
}

func createDevice(db *gorm.DB, userID uint, body *addDevice) (*model.Device, error) {
	device := model.Device{
		UserID:         userID,
		RegistrationID: body.RegistrationID,
		Name:           body.Name,
		Type:           body.Type.String(),
		State:          types.DeviceStateLinked.String(),
	}
	if err := db.Create(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func createNotificationKey(db *gorm.DB, gc *util.WebClient, userID uint, registrationID string) (*model.NotificationKey, error) {
	var notificationKey model.NotificationKey

	// If no notification key exists: create and register with GCM
	if db.Where(model.NotificationKey{UserID: userID}).First(&notificationKey).RecordNotFound() {
		bytes := make([]byte, 48)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}
		groupName := hex.EncodeToString(bytes)
		key, err := util.CreateNotificationGroup(gc, groupName, registrationID)
		if err != nil {
			return nil, err
		}

		notificationKey = model.NotificationKey{
			UserID:    userID,
			Key:       key,
			GroupName: groupName,
		}

		if err := db.Create(&notificationKey).Error; err != nil {
			return nil, err
		}
		return &notificationKey, nil
	}

	// If notification key exists: add device to notification group
	err := util.AddNotificationGroup(gc, notificationKey.GroupName, notificationKey.Key, registrationID)
	if err != nil {
		return nil, err
	}
	return &notificationKey, nil
}

func createEncryptionKey(db *gorm.DB, userID uint) (*model.EncryptionKey, error) {
	var encryptionKey model.EncryptionKey
	if db.Where(model.EncryptionKey{UserID: userID}).First(&encryptionKey).RecordNotFound() {
		key := make([]byte, 48)
		if _, err := rand.Read(key); err != nil {
			return nil, err
		}
		encryptionKey = model.EncryptionKey{
			UserID: userID,
			Key:    hex.EncodeToString(key),
		}
		if err := db.Create(&encryptionKey).Error; err != nil {
			return nil, err
		}
	}
	return &encryptionKey, nil
}
