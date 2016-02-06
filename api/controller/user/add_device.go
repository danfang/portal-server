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
	"portal-server/api/controller/context"
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

// AddDeviceEndpoint allows users to register new GCM devices, which returns encryption
// and notification keys on success.
func AddDeviceEndpoint(c *gin.Context) {
	var body addDevice
	if !controller.ValidJSON(c, &body) {
		return
	}

	user := context.UserFromContext(c)
	wc := context.WebClientFromContext(c, gcmEndpoint)

	s := context.StoreFromContext(c)
	s.Transaction(func(store store.Store) error {
		var err error
		if store.Devices().DeviceCount(&model.Device{RegistrationID: body.RegistrationID}) >= 1 {
			err = errs.ErrDuplicateDeviceToken
			c.JSON(http.StatusBadRequest, controller.RenderError(err))
			return err
		}

		device, err := createDevice(store, user, &body)
		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}

		notificationKey, err := createNotificationKey(store, wc, user, device.RegistrationID)
		if err, isGCMError := err.(errs.GCMError); isGCMError {
			c.JSON(http.StatusBadRequest, controller.DetailError{
				Error:  errs.ErrUnableToRegisterDevice.Error(),
				Reason: err.Error(),
			})
			return err
		}

		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}

		encryptionKey, err := getEncryptionKey(store, user)
		if err != nil {
			controller.InternalServiceError(c, err)
			return err
		}
		c.JSON(http.StatusOK, addDeviceResponse{
			EncryptionKey:   encryptionKey.Key,
			NotificationKey: notificationKey.Key,
		})
		return nil
	})
}

func createDevice(store store.Store, user *model.User, body *addDevice) (*model.Device, error) {
	device := &model.Device{
		User:           *user,
		RegistrationID: body.RegistrationID,
		Name:           body.Name,
		Type:           body.Type,
		State:          model.DeviceStateLinked,
	}
	if err := store.Devices().CreateDevice(device); err != nil {
		return nil, err
	}
	return device, nil
}

func createNotificationKey(store store.Store, wc *util.WebClient, user *model.User, registrationID string) (*model.NotificationKey, error) {
	notificationKey, found := store.NotificationKeys().FindKey(&model.NotificationKey{
		UserID: user.ID,
	})

	// If no notification key exists: create and register with GCM
	if !found {
		bytes := make([]byte, 48)
		if _, err := rand.Read(bytes); err != nil {
			return nil, err
		}

		groupName := hex.EncodeToString(bytes)
		key, err := util.CreateNotificationGroup(wc, groupName, registrationID)
		if err != nil {
			return nil, err
		}

		notificationKey = &model.NotificationKey{
			User:      *user,
			Key:       key,
			GroupName: groupName,
		}

		if err := store.NotificationKeys().CreateKey(notificationKey); err != nil {
			return nil, err
		}
		return notificationKey, nil
	}

	// If notification key exists: add device to notification group
	err := util.AddNotificationGroup(wc, notificationKey.GroupName, notificationKey.Key, registrationID)
	if err != nil {
		return nil, err
	}

	return notificationKey, nil
}

func getEncryptionKey(store store.Store, user *model.User) (*model.EncryptionKey, error) {
	encryptionKey, found := store.EncryptionKeys().FindKey(&model.EncryptionKey{
		User: *user,
	})
	// Create new key if not found
	if !found {
		key := make([]byte, 32)
		if _, err := rand.Read(key); err != nil {
			return nil, err
		}
		encryptionKey = &model.EncryptionKey{
			User: *user,
			Key:  hex.EncodeToString(key),
		}
		if err := store.EncryptionKeys().CreateKey(encryptionKey); err != nil {
			return nil, err
		}
	}
	return encryptionKey, nil
}
