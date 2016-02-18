package user

import (
	"bytes"
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
	"github.com/stretchr/testify/assert"
)

func TestSignout(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)
	g.Describe("GET /user/signout", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownTestStore(s)
		})

		g.It("Should unregister and delete the user device and user token", func() {
			user := model.User{
				Email: "test@portal.com",
				UUID:  "1",
			}
			s.Users().CreateUser(&user)
			userToken := model.UserToken{
				User:  user,
				Token: "token",
			}
			s.UserTokens().CreateToken(&userToken)
			key := model.NotificationKey{
				User:      user,
				Key:       "key",
				GroupName: "name",
			}
			s.NotificationKeys().CreateKey(&key)
			device := model.Device{
				User:            user,
				NotificationKey: key,
				UUID:            "1",
				Name:            "My Device",
				Type:            "phone",
				State:           model.DeviceStateLinked,
				RegistrationID:  "registration_id",
			}
			s.Devices().CreateDevice(&device)
			w := testSignout(s, &user, &userToken, signout{
				DeviceID: device.UUID,
			})
			assert.Equal(t, 200, w.Code)

			fromDB, _ := s.Devices().FindDevice(&model.Device{UserID: user.ID})
			assert.Equal(t, model.DeviceStateUnlinked, fromDB.State)

			_, found := s.UserTokens().FindToken(&model.UserToken{})
			assert.False(t, found)
		})
	})
}

func testSignout(s store.Store, user *model.User, userToken *model.UserToken, input interface{}) *httptest.ResponseRecorder {
	r := testutil.TestRouter(middleware.SetStore(s))

	// Set the user and token
	r.Use(func(c *gin.Context) {
		context.UserToContext(c, user)
		context.UserTokenToContext(c, userToken)
		c.Next()
	})

	r.POST("/", SignoutEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(body)))
	r.ServeHTTP(w, req)
	return w
}
