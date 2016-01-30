package user

import (
	"encoding/json"
	"portal-server/model"
	"portal-server/model/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var getDevicesDB gorm.DB

func init() {
	gin.SetMode(gin.TestMode)
	getDevicesDB, _ = gorm.Open("sqlite3", ":memory:")
	getDevicesDB.CreateTable(&model.User{}, &model.Device{})
}

func TestGetDevicesEndpoint_NoDevices(t *testing.T) {
	w := testGetDevices(404)
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
	getDevicesDB.Create(&user)
	getDevicesDB.Create(&model.Device{
		User:           user,
		Name:           "Nexus 6P",
		Type:           "phone",
		RegistrationID: "1",
		State:          types.DeviceStateLinked.String(),
	})
	getDevicesDB.Create(&model.Device{
		User:           user,
		Name:           "Chrome 4.2",
		Type:           "chrome",
		RegistrationID: "2",
		State:          types.DeviceStateLinked.String(),
	})
	getDevicesDB.Create(&model.Device{
		User:           user,
		Name:           "Unlinked Desktop",
		Type:           "desktop",
		RegistrationID: "3",
		State:          types.DeviceStateUnlinked.String(),
	})
	w := testGetDevices(user.ID)
	assert.Equal(t, 200, w.Code)

	var res deviceListResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, 2, len(res.Devices))
}

func testGetDevices(userID uint) *httptest.ResponseRecorder {
	// Create the router
	userRouter := Router{&getDevicesDB, http.DefaultClient}
	r := gin.New()

	// Set the userID
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	r.GET("/", userRouter.GetDevicesEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}
