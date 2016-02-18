package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/api/errs"
	"portal-server/model"

	"github.com/gin-gonic/gin"
)

func DeleteMessageEndpoint(c *gin.Context) {
	user := context.UserFromContext(c)
	store := context.StoreFromContext(c)
	messageID := c.Param("mid")

	rows := store.Messages().DeleteMessages(&model.Message{
		UserID:    user.ID,
		MessageID: messageID,
	})
	if rows == 0 {
		c.JSON(http.StatusNotFound, controller.RenderError(errs.ErrMessageNotFound))
		return
	}
	c.JSON(http.StatusOK, controller.RenderSuccess(true))
}
