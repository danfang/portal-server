package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"portal-server/api/util"
	"portal-server/model"
	"testing"

	"github.com/franela/goblin"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"portal-server/store"
)

var s = store.GetTestStore()

func TestAddDevice(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)
	g.Describe("Data store functions", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownTestStore(s)
		})

		g.It("Should successfully create a device without duplicates", func() {
			user := &model.User{
				Email:    "test@portal.com",
				Verified: true,
			}

			s.Users().CreateUser(user)

			body := &addDevice{
				RegistrationID: "a_token",
				Type:           model.DeviceTypePhone,
			}

			device, err := createDevice(s, user, body, &model.NotificationKey{})
			assert.NoError(t, err)
			assert.Equal(t, "a_token", device.RegistrationID)
			assert.Equal(t, model.DeviceTypePhone, device.Type)
			assert.Equal(t, model.DeviceStateLinked, device.State)
			_, err = uuid.FromString(device.UUID)
			assert.NoError(t, err)

			fromDB, _ := s.Devices().GetRelatedUser(device)
			assert.Equal(t, user.ID, fromDB.ID)
			assert.Equal(t, user.Email, fromDB.Email)
			assert.Equal(t, user.Verified, fromDB.Verified)

			// Cannot create duplicate devices for the same user
			_, err = createDevice(s, user, body, &model.NotificationKey{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "UNIQUE")

			// Cannot create duplicate devices for other users
			otherUser := &model.User{
				Email: "abc@def.com",
			}
			_, err = createDevice(s, otherUser, body, &model.NotificationKey{})
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "UNIQUE")
		})

		g.It("Should successfully create a notification key from a GCM response", func() {
			user := &model.User{Email: "test2@portal.com"}
			err := s.Users().CreateUser(user)
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
			key, err := createNotificationKey(s, client, user, "registrationId")
			assert.NoError(t, err)
			assert.Regexp(t, "^[a-fA-F0-9]+$", key.GroupName)
			assert.Equal(t, notificationKey, key.Key)

			fromDB, _ := s.NotificationKeys().GetRelatedUser(key)
			assert.Equal(t, user.ID, fromDB.ID)
			assert.Equal(t, user.Email, fromDB.Email)
		})

		g.It("Should add a device to a notification group if one already exists", func() {
			user := &model.User{Email: "test3@portal.com"}
			s.Users().CreateUser(user)

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
			key1, err := createNotificationKey(s, client, user, "registrationId")
			assert.NoError(t, err)
			assert.Equal(t, notificationKey, key1.Key)
			server.Close()

			// Attempt to create the second key, expecting an "add" operation
			requestTest = func(r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)
				assert.Contains(t, string(body), "add")
				assert.Contains(t, string(body), notificationKey)
			}
			server, client = util.TestHTTP(requestTest, 200, string(mockResponse))
			key2, err := createNotificationKey(s, client, user, "registrationId")

			// Make sure the keys are the same
			assert.NoError(t, err)
			assert.Equal(t, key1.Key, key2.Key)
			assert.Equal(t, key1.GroupName, key2.GroupName)

			// Make sure only one key exists in the DB
			count := s.NotificationKeys().GetCount(&model.NotificationKey{UserID: user.ID})
			assert.Equal(t, 1, count)

			server.Close()
		})

		g.It("Should create unique encryption keys per user", func() {
			user := &model.User{Email: "test4@portal.com"}
			s.Users().CreateUser(user)

			// Create the key
			key1, err := getEncryptionKey(s, user)
			assert.NoError(t, err)
			assert.Regexp(t, "^[a-fA-F0-9]{64}$", key1.Key)

			fromDB, _ := s.EncryptionKeys().GetRelatedUser(key1)
			assert.Equal(t, user.ID, fromDB.ID)
			assert.Equal(t, user.Email, fromDB.Email)

			// Attempt to create a duplicate key
			key2, err := getEncryptionKey(s, user)
			assert.NoError(t, err)
			assert.Equal(t, key1.Key, key2.Key)

			// Make sure only one key exists in the DB
			count := s.EncryptionKeys().GetCount(&model.EncryptionKey{UserID: user.ID})
			assert.Equal(t, 1, count)
		})
	})
}
