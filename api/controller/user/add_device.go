package user

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/api/util"
	"portal-server/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"portal-server/store"
)

const gcmEndpoint = "https://android.googleapis.com/gcm/notification"

type addDevice struct {
	RegistrationID string `json:"registration_id" valid:"required"`
	Name           string `json:"name" valid:"required"`
	Type           string `json:"type" valid:"required,matches(phone,chrome,desktop)"`
}

type addDeviceResponse struct {
	EncryptionKey   string `json:"encryption_key"`
	NotificationKey string `json:"notification_key"`
}

func (r Router) AddDeviceEndpoint(c *gin.Context) {
	user := c.MustGet("userID").(*model.User)

	var body addDevice
	if !controller.ValidJSON(c, &body) {
		return
	}

	var device *model.Device
	r.Store.Transaction(func(tx store.Store) error {
		var err error
		if tx.Devices().DeviceCount(&model.Device{RegistrationID: body.RegistrationID}) >= 1 {
			err = errs.ErrDuplicateDeviceToken
			c.JSON(http.StatusBadRequest, controller.RenderError(err))
			return err
		}

		device, err = createDevice(tx, user, &body)
		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}
		return nil
	})

	gcmClient := &util.WebClient{BaseURL: gcmEndpoint, HTTPClient: r.HTTPClient}

	notificationKey, err := createNotificationKey(tx, gcmClient, user, device.RegistrationID)
	if err, isGCMError := err.(errs.GCMError); isGCMError {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, controller.DetailError{
			Error:  errs.ErrUnableToRegisterDevice.Error(),
			Reason: err.Error(),
		})
		return
	}
	if err != nil {
		tx.Rollback()
		controller.InternalServiceError(c, err)
		return
	}

	encryptionKey, err := getEncryptionKey(tx, user)
	if err != nil {
		tx.Rollback()
		controller.InternalServiceError(c, err)
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, addDeviceResponse{
		EncryptionKey:   encryptionKey.Key,
		NotificationKey: notificationKey.Key,
	})
}

func createDevice(db *gorm.DB, user *model.User, body *addDevice) (*model.Device, error) {
	device := model.Device{
		User:           *user,
		RegistrationID: body.RegistrationID,
		Name:           body.Name,
		Type:           body.Type,
		State:          model.DeviceStateLinked,
	}
	if err := db.Create(&device).Error; err != nil {
		return nil, err
	}
	return &device, nil
}

func createNotificationKey(db *gorm.DB, gc *util.WebClient, user *model.User, registrationID string) (*model.NotificationKey, error) {
	var notificationKey model.NotificationKey

	// If no notification key exists: create and register with GCM
	if db.Where(model.NotificationKey{User: *user}).First(&notificationKey).RecordNotFound() {
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
			User:      *user,
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

func getEncryptionKey(db *gorm.DB, user *model.User) (*model.EncryptionKey, error) {
	var encryptionKey model.EncryptionKey
	if db.Where(model.EncryptionKey{User: *user}).First(&encryptionKey).RecordNotFound() {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, err
		}
		encryptionKey = model.EncryptionKey{
			User: *user,
			Key:  hex.EncodeToString(key),
		}
		if err := db.Create(&encryptionKey).Error; err != nil {
			return nil, err
		}
	}
	return &encryptionKey, nil
}
