package user

import (
	"portal-server/api/routing"
	"portal-server/model"
	"portal-server/model/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// DeviceListResponse contains existing, connected user devices.
//
// swagger:response deviceList
type DeviceListResponse struct {
	// in: body
	Body deviceListResponse `json:"device_list"`
}

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
	userID := c.MustGet("userID").(uint)
	devices, err := getLinkedDevices(r.Db, userID)
	if err != nil {
		routing.InternalServiceError(c, err)
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

func getLinkedDevices(db *gorm.DB, userID uint) ([]model.Device, error) {
	var devices []model.Device
	if err := db.Where(model.Device{
		UserID: userID,
		State:  types.DeviceStateLinked.String(),
	}).Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}
