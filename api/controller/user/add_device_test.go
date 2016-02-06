package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"portal-server/api/util"
	"portal-server/model"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"portal-server/store"
)

var addDeviceStore = store.GetTestStore()

func TestCreateDevice(t *testing.T) {
	user := &model.User{
		Email:    "test@portal.com",
		Verified: true,
	}

	addDeviceStore.Users().CreateUser(user)

	body := &addDevice{
		RegistrationID: "a_token",
		Type:           model.DeviceTypePhone,
	}

	device, err := createDevice(addDeviceStore, user, body)
	assert.NoError(t, err)
	assert.Equal(t, "a_token", device.RegistrationID)
	assert.Equal(t, model.DeviceTypePhone, device.Type)
	assert.Equal(t, model.DeviceStateLinked, device.State)

	fromDB, _ := addDeviceStore.Devices().GetRelatedUser(device)
	assert.Equal(t, user.ID, fromDB.ID)
	assert.Equal(t, user.Email, fromDB.Email)
	assert.Equal(t, user.Verified, fromDB.Verified)

	// Cannot create duplicate devices for the same user
	_, err = createDevice(addDeviceStore, user, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")

	// Cannot create duplicate devices for other users
	otherUser := &model.User{
		Email: "abc@def.com",
	}
	_, err = createDevice(addDeviceStore, otherUser, body)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "UNIQUE")
}

func TestCreateNotificationKey(t *testing.T) {
	user := &model.User{Email: "test2@portal.com"}
	err := addDeviceStore.Users().CreateUser(user)
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
	key, err := createNotificationKey(addDeviceStore, client, user, "registrationId")
	assert.NoError(t, err)
	assert.Regexp(t, "^[a-fA-F0-9]+$", key.GroupName)
	assert.Equal(t, notificationKey, key.Key)

	fromDB, _ := addDeviceStore.NotificationKeys().GetRelatedUser(key)
	assert.Equal(t, user.ID, fromDB.ID)
	assert.Equal(t, user.Email, fromDB.Email)
}

func TestCreateNotificationKey_Duplicate(t *testing.T) {
	user := &model.User{Email: "test3@portal.com"}
	addDeviceStore.Users().CreateUser(user)

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
	key1, err := createNotificationKey(addDeviceStore, client, user, "registrationId")
	assert.NoError(t, err)
	assert.Equal(t, notificationKey, key1.Key)
	server.Close()

	// Attempt to create the second key, expecting an "add" operation
	requestTest = func(r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		assert.Contains(t, string(body), "add")
	}
	server, client = util.TestHTTP(requestTest, 200, string(mockResponse))
	key2, err := createNotificationKey(addDeviceStore, client, user, "registrationId")

	// Make sure the keys are the same
	assert.NoError(t, err)
	assert.Equal(t, key1.Key, key2.Key)
	assert.Equal(t, key1.GroupName, key2.GroupName)

	// Make sure only one key exists in the DB
	count := addDeviceStore.NotificationKeys().GetCount(&model.NotificationKey{UserID: user.ID})
	assert.Equal(t, 1, count)

	server.Close()
}

func TestGetEncryptionKey(t *testing.T) {
	user := &model.User{Email: "test4@portal.com"}
	addDeviceStore.Users().CreateUser(user)

	// Create the key
	key1, err := getEncryptionKey(addDeviceStore, user)
	assert.NoError(t, err)
	assert.Regexp(t, "^[a-fA-F0-9]{64}$", key1.Key)

	fromDB, _ := addDeviceStore.EncryptionKeys().GetRelatedUser(key1)
	assert.Equal(t, user.ID, fromDB.ID)
	assert.Equal(t, user.Email, fromDB.Email)

	// Attempt to create a duplicate key
	key2, err := getEncryptionKey(addDeviceStore, user)
	assert.NoError(t, err)
	assert.Equal(t, key1.Key, key2.Key)

	// Make sure only one key exists in the DB
	count := addDeviceStore.EncryptionKeys().GetCount(&model.EncryptionKey{UserID: user.ID})
	assert.Equal(t, 1, count)
}
