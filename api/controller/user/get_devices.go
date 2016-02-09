package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
)

type deviceListResponse struct {
	Devices []linkedDevice `json:"devices"`
}

type linkedDevice struct {
	DeviceID  string `json:"device_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

// GetDevicesEndpoint retrieves connected user devices.
func GetDevicesEndpoint(c *gin.Context) {
	user := context.UserFromContext(c)
	store := context.StoreFromContext(c)
	devices, err := store.Devices().GetAllLinkedDevices(user)
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}
	linkedDevices := make([]linkedDevice, 0, len(devices))
	for _, value := range devices {
		linkedDevices = append(linkedDevices, linkedDevice{
			DeviceID:  value.UUID,
			CreatedAt: value.CreatedAt.Unix(),
			UpdatedAt: value.UpdatedAt.Unix(),
			Name:      value.Name,
			Type:      value.Type,
		})
	}
	c.JSON(http.StatusOK, deviceListResponse{
		Devices: linkedDevices,
	})
}
