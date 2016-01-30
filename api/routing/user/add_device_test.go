package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"portal-server/api/util"
	"portal-server/model"
	"portal-server/model/types"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

var addDeviceDB gorm.DB

func init() {
	addDeviceDB, _ = gorm.Open("sqlite3", ":memory:")
	addDeviceDB.LogMode(false)
	addDeviceDB.CreateTable(&model.User{}, &model.EncryptionKey{}, &model.Device{}, &model.NotificationKey{})
}

func TestCreateDevice(t *testing.T) {
	user := model.User{
		Email:    "test@portal.com",
		Verified: true,
	}

	addDeviceDB.Create(&user)

	body := &addDevice{
		RegistrationID: "a_token",
		Type:           types.DeviceTypePhone,
	}

	device, err := createDevice(&addDeviceDB, user.ID, body)
	assert.NoError(t, err)
	assert.Equal(t, "a_token", device.RegistrationID)
	assert.Equal(t, types.DeviceTypePhone.String(), device.Type)
	assert.Equal(t, types.DeviceStateLinked.String(), device.State)

	var fromDB model.User
	addDeviceDB.Model(device).Related(&fromDB)
	assert.Equal(t, user.ID, fromDB.ID)
	assert.Equal(t, user.Email, fromDB.Email)
	assert.Equal(t, user.Verified, fromDB.Verified)

	// Cannot create duplicate devices for the same user
	_, err = createDevice(&addDeviceDB, user.ID, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")

	// Cannot create duplicate devices for other users
	otherUserID := uint(5)
	_, err = createDevice(&addDeviceDB, otherUserID, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestCreateNotificationKey(t *testing.T) {
	user := model.User{Email: "test2@portal.com"}
	addDeviceDB.Create(&user)
	notificationKey := "notificationKey"
	mockResponse, _ := json.Marshal(map[string]string{
		"notification_key": notificationKey,
	})
	requestTest := func(r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		assert.Contains(t, string(body), "create")
	}
	server, client := util.TestHTTP(requestTest, 200, string(mockResponse))
	defer server.Close()
	key, err := createNotificationKey(&addDeviceDB, client, user.ID, "registrationId")
	assert.NoError(t, err)
	assert.Regexp(t, "^[a-fA-F0-9]+$", key.GroupName)
	assert.Equal(t, notificationKey, key.Key)

	var fromDB model.User
	addDeviceDB.Model(key).Related(&fromDB)
	assert.Equal(t, user.ID, fromDB.ID)
	assert.Equal(t, user.Email, fromDB.Email)
}

func TestCreateNotificationKey_Duplicate(t *testing.T) {
	user := model.User{Email: "test3@portal.com"}
	addDeviceDB.Create(&user)

	notificationKey := "notificationKey"
	mockResponse, _ := json.Marshal(map[string]string{
		"notification_key": notificationKey,
	})

	// Create the first key, expecting a "create" operation
	requestTest := func(r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		assert.Contains(t, string(body), "create")
	}
	server, client := util.TestHTTP(requestTest, 200, string(mockResponse))
	key1, err := createNotificationKey(&addDeviceDB, client, user.ID, "registrationId")
	assert.NoError(t, err)
	assert.Equal(t, notificationKey, key1.Key)
	server.Close()

	// Attempt to create the second key, expecting an "add" operation
	requestTest = func(r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		assert.Contains(t, string(body), "add")
	}
	server, client = util.TestHTTP(requestTest, 200, string(mockResponse))
	key2, err := createNotificationKey(&addDeviceDB, client, user.ID, "registrationId")

	// Make sure the keys are the same
	assert.NoError(t, err)
	assert.Equal(t, key1.Key, key2.Key)
	assert.Equal(t, key1.GroupName, key2.GroupName)

	// Make sure only one key exists in the DB
	var count int
	addDeviceDB.Model(&model.NotificationKey{}).Where("user_id = ?", user.ID).Count(&count)
	assert.Equal(t, 1, count)

	server.Close()
}

func TestCreateEncryptionKey(t *testing.T) {
	user := model.User{Email: "test4@portal.com"}
	addDeviceDB.Create(&user)

	// Create the key
	key1, err := createEncryptionKey(&addDeviceDB, user.ID)
	assert.NoError(t, err)
	assert.Regexp(t, "^[a-fA-F0-9]+$", key1.Key)

	var fromDB model.User
	addDeviceDB.Model(key1).Related(&fromDB)
	assert.Equal(t, user.ID, fromDB.ID)
	assert.Equal(t, user.Email, fromDB.Email)

	// Attempt to create a duplicate key
	key2, err := createEncryptionKey(&addDeviceDB, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, key1.Key, key2.Key)

	// Make sure only one key exists in the DB
	var count int
	addDeviceDB.Model(&model.EncryptionKey{}).Where("user_id = ?", user.ID).Count(&count)
	assert.Equal(t, 1, count)
}
