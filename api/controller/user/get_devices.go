package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/model"

	"github.com/gin-gonic/gin"
)

type deviceListResponse struct {
	Devices []linkedDevice `json:"devices"`
}

type linkedDevice struct {
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	Name      string `json:"name"`
	Type      string `json:"type"`
}

// GetDevicesEndpoint retrieves connected user devices.
func (r Router) GetDevicesEndpoint(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	devices, err := r.Store.Devices().GetAllLinkedDevices(user)
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}
	linkedDevices := make([]linkedDevice, 0, len(devices))
	for _, value := range devices {
		linkedDevices = append(linkedDevices, linkedDevice{
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
