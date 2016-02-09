package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/model"
	"testing"

	"bytes"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"portal-server/api/controller/context"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/api/util"
	"portal-server/store"
	"portal-server/vendor/github.com/franela/goblin"
)

func TestSignout(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)
	g.Describe("GET /user/signout", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownStoreForTest(s)
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
			requestTest := func(r *http.Request) {
				body, _ := ioutil.ReadAll(r.Body)
				assert.Contains(t, string(body), "remove")
				assert.Contains(t, string(body), "registration_id")
			}
			w := testSignout(requestTest, s, &user, &userToken, signout{
				DeviceID: device.UUID,
			})
			assert.Equal(t, 200, w.Code)

			count := s.Devices().DeviceCount(&model.Device{})
			assert.Equal(t, 0, count)

			_, found := s.UserTokens().FindToken(&model.UserToken{})
			assert.False(t, found)
		})
	})
}

func testSignout(requestTest func(*http.Request),
	s store.Store, user *model.User, userToken *model.UserToken, input interface{}) *httptest.ResponseRecorder {

	server, client := util.TestHTTP(requestTest, 200, "")

	gcmEndpoint = server.URL

	r := testutil.TestRouter(
		middleware.SetStore(s),
		middleware.SetWebClient(client.HTTPClient),
	)

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
