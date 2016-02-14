package store

import (
	"portal-server/model"
	"testing"

	"github.com/franela/goblin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestDeviceStore(t *testing.T) {
	var db *gorm.DB
	var store deviceStore
	var user model.User
	var notificationKey model.NotificationKey
	g := goblin.Goblin(t)

	g.Describe("DeviceStore", func() {
		g.BeforeEach(func() {
			db = GetTestDB()
			user = model.User{
				UUID:  "1",
				Email: "test@portal.com",
			}
			db.Create(&user)
			notificationKey = model.NotificationKey{
				User:      user,
				GroupName: "group",
				Key:       "key",
			}
			db.Create(&notificationKey)
			store = deviceStore{db}
		})

		g.AfterEach(func() {
			TeardownTestDB(db)
		})

		g.It("CreateDevice", func() {
			store.CreateDevice(&model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "1",
				Type:            model.DeviceTypePhone,
				State:           model.DeviceStateLinked,
			})
			var device model.Device
			db.Where(&model.Device{UserID: user.ID}).First(&device)
			assert.Equal(t, notificationKey.ID, device.NotificationKeyID)
			assert.Equal(t, "1", device.UUID)
			assert.Equal(t, model.DeviceTypePhone, device.Type)
			assert.Equal(t, model.DeviceStateLinked, device.State)
		})

		g.It("SaveDevice", func() {
			device := model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "1",
				Type:            model.DeviceTypePhone,
				State:           model.DeviceStateLinked,
			}
			db.Create(&device)
			device.UUID = "2"
			device.State = model.DeviceStateUnlinked
			device.Type = model.DeviceTypeChrome
			store.SaveDevice(&device)

			var d model.Device
			db.Where(&model.Device{UserID: user.ID}).First(&d)
			assert.Equal(t, notificationKey.ID, d.NotificationKeyID)
			assert.Equal(t, "2", d.UUID)
			assert.Equal(t, model.DeviceTypeChrome, d.Type)
			assert.Equal(t, model.DeviceStateUnlinked, d.State)
		})

		g.It("FindDevice", func() {
			db.Create(&model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "1",
				Name:            "Chrome OSX",
				Type:            model.DeviceTypePhone,
				State:           model.DeviceStateLinked,
			})
			device, found := store.FindDevice(&model.Device{UserID: user.ID})
			assert.True(t, found)
			assert.Equal(t, user.ID, device.UserID)
			assert.Equal(t, "Chrome OSX", device.Name)
			assert.Equal(t, model.DeviceTypePhone, device.Type)
			assert.Equal(t, model.DeviceStateLinked, device.State)
		})

		g.It("DeleteDevice", func() {
			device := &model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "1",
				Type:            model.DeviceTypePhone,
				State:           model.DeviceStateLinked,
			}
			db.Create(device)
			assert.NoError(t, store.DeleteDevice(device))
			assert.True(t, db.Where(&model.Device{UserID: user.ID}).First(&model.Device{}).RecordNotFound())
		})

		g.It("DeviceCount", func() {
			count := 10
			for i := 0; i < count; i++ {
				db.Create(&model.Device{
					User:            user,
					NotificationKey: notificationKey,
					UUID:            string(i),
					RegistrationID:  string(i),
					Type:            model.DeviceTypePhone,
					State:           model.DeviceStateLinked,
				})
			}
			assert.Equal(t, count, store.DeviceCount(&model.Device{UserID: user.ID}))
		})

		g.It("GetAllLinkedDevices", func() {
			db.Create(&model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "1",
				RegistrationID:  "1",
				Type:            model.DeviceTypePhone,
				State:           model.DeviceStateLinked,
			})
			db.Create(&model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "2",
				RegistrationID:  "2",
				Type:            model.DeviceTypeChrome,
				State:           model.DeviceStateUnlinked,
			})
			db.Create(&model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "3",
				RegistrationID:  "3",
				Type:            model.DeviceTypeDesktop,
				State:           model.DeviceStateUnlinked,
			})
			devices, err := store.GetAllLinkedDevices(&user)
			assert.NoError(t, err)
			assert.Equal(t, 1, len(devices))
			assert.Equal(t, "1", devices[0].UUID)
			assert.Equal(t, "1", devices[0].RegistrationID)
			assert.Equal(t, model.DeviceTypePhone, devices[0].Type)
			assert.Equal(t, model.DeviceStateLinked, devices[0].State)
		})

		g.It("GetRelatedUser", func() {
			device := &model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "1",
				Type:            model.DeviceTypePhone,
				State:           model.DeviceStateLinked,
			}
			db.Create(device)

			u, err := store.GetRelatedUser(device)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, u.ID)
			assert.Equal(t, user.Email, u.Email)
		})

		g.It("GetRelatedKey", func() {
			device := &model.Device{
				User:            user,
				NotificationKey: notificationKey,
				UUID:            "1",
				Type:            model.DeviceTypePhone,
				State:           model.DeviceStateLinked,
			}
			db.Create(device)

			key, err := store.GetRelatedKey(device)
			assert.NoError(t, err)
			assert.Equal(t, notificationKey.ID, key.ID)
			assert.Equal(t, notificationKey.GroupName, key.GroupName)
			assert.Equal(t, notificationKey.Key, key.Key)
		})
	})
}
