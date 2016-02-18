package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/model"

	"github.com/gin-gonic/gin"
)

type signout struct {
	DeviceID string `json:"device_id"`
}

func SignoutEndpoint(c *gin.Context) {
	s := context.StoreFromContext(c)

	// Delete the user token
	userToken := context.UserTokenFromContext(c)
	err := s.UserTokens().DeleteToken(userToken)
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}

	// Unlink the device, if provided
	var body signout
	c.BindJSON(&body)
	if body.DeviceID != "" {
		user := context.UserFromContext(c)
		device, found := s.Devices().FindDevice(&model.Device{UserID: user.ID, UUID: body.DeviceID})
		if found {
			device.State = model.DeviceStateUnlinked
			s.Devices().SaveDevice(device)
		}
	}

	c.JSON(http.StatusOK, controller.RenderSuccess(true))
}
