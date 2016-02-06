package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/model"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"portal-server/store"
)

var getDevicesStore = store.GetTestStore()

func init() {
	gin.SetMode(gin.TestMode)
}

func TestGetDevicesEndpoint_NoDevices(t *testing.T) {
	w := testGetDevices(&model.User{})
	assert.Equal(t, 200, w.Code)
	expected, _ := json.Marshal(gin.H{
		"devices": []linkedDevice{},
	})
	assert.JSONEq(t, string(expected), w.Body.String())
}

func TestGetDevicesEndpoint_LinkedDevices(t *testing.T) {
	user := model.User{
		Email: "test@portal.com",
	}
	getDevicesStore.Users().CreateUser(&user)
	getDevicesStore.Devices().CreateDevice(&model.Device{
		User:           user,
		Name:           "Nexus 6P",
		Type:           "phone",
		RegistrationID: "1",
		State:          model.DeviceStateLinked,
	})
	getDevicesStore.Devices().CreateDevice(&model.Device{
		User:           user,
		Name:           "Chrome 4.2",
		Type:           "chrome",
		RegistrationID: "2",
		State:          model.DeviceStateLinked,
	})
	getDevicesStore.Devices().CreateDevice(&model.Device{
		User:           user,
		Name:           "Unlinked Desktop",
		Type:           "desktop",
		RegistrationID: "3",
		State:          model.DeviceStateUnlinked,
	})
	w := testGetDevices(&user)
	assert.Equal(t, 200, w.Code)

	var res deviceListResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, 2, len(res.Devices))
}

func testGetDevices(user *model.User) *httptest.ResponseRecorder {
	// Create the router
	userRouter := Router{
		Store:      getDevicesStore,
		HTTPClient: http.DefaultClient,
	}
	r := gin.New()

	// Set the userID
	r.Use(func(c *gin.Context) {
		c.Set("user", user)
		c.Next()
	})

	r.GET("/", userRouter.GetDevicesEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}
