package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/api/controller/context"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/model"
	"portal-server/store"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetDevices(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)
	g.Describe("GET /user/devices", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownTestStore(s)
		})

		g.It("Should return an empty array for a user with no linked devices", func() {
			w := testGetDevices(s, &model.User{})
			assert.Equal(t, 200, w.Code)
			expected, _ := json.Marshal(gin.H{
				"devices": []linkedDevice{},
			})
			assert.JSONEq(t, string(expected), w.Body.String())
		})

		g.It("Should retrieve all linked devices for a user", func() {
			user := model.User{Email: "test@portal.com"}
			s.Users().CreateUser(&user)
			key := model.NotificationKey{
				User:      user,
				Key:       "key",
				GroupName: "name",
			}
			s.NotificationKeys().CreateKey(&key)
			s.Devices().CreateDevice(&model.Device{
				User:            user,
				NotificationKey: key,
				UUID:            uuid.NewV4().String(),
				Name:            "Nexus 6P",
				Type:            "phone",
				RegistrationID:  "1",
				State:           model.DeviceStateLinked,
			})
			s.Devices().CreateDevice(&model.Device{
				User:            user,
				NotificationKey: key,
				UUID:            uuid.NewV4().String(),
				Name:            "Chrome 4.2",
				Type:            "chrome",
				RegistrationID:  "2",
				State:           model.DeviceStateLinked,
			})
			s.Devices().CreateDevice(&model.Device{
				User:            user,
				NotificationKey: key,
				UUID:            uuid.NewV4().String(),
				Name:            "Unlinked Desktop",
				Type:            "desktop",
				RegistrationID:  "3",
				State:           model.DeviceStateUnlinked,
			})
			w := testGetDevices(s, &user)
			assert.Equal(t, 200, w.Code)

			var res deviceListResponse
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Equal(t, 2, len(res.Devices))
		})
	})
}

func testGetDevices(s store.Store, user *model.User) *httptest.ResponseRecorder {
	r := testutil.TestRouter(middleware.SetStore(s))

	// Set the userID
	r.Use(func(c *gin.Context) {
		context.UserToContext(c, user)
		c.Next()
	})

	r.GET("/", GetDevicesEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}
