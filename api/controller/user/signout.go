package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/api/util"
	"portal-server/model"
)

type signout struct {
	DeviceID string `json:"device_id"`
}

func SignoutEndpoint(c *gin.Context) {
	var body signout
	if !controller.ValidJSON(c, &body) {
		return
	}

	user := context.UserFromContext(c)
	s := context.StoreFromContext(c)

	// Get the device and associated notification key
	device, found := s.Devices().FindDevice(&model.Device{UserID: user.ID, UUID: body.DeviceID})
	if found {
		key, err := s.Devices().GetRelatedKey(device)
		if err == nil {
			// Unregister the key
			wc := context.WebClientFromContext(c, gcmEndpoint)
			util.RemoveNotificationGroup(wc, key.GroupName, key.Key, device.RegistrationID)

			// Delete the device
			s.Devices().DeleteDevice(device)
		}
	}

	// Delete the user token
	userToken := context.UserTokenFromContext(c)
	err := s.UserTokens().DeleteToken(userToken)
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, controller.RenderSuccess())
}
