package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"

	"portal-server/api/errs"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func SyncMessagesEndpoint(c *gin.Context) {
	user := context.UserFromContext(c)
	store := context.StoreFromContext(c)
	messageID := c.Param("mid")
	messages, err := store.Messages().GetMessagesSince(user, messageID)
	if err == gorm.RecordNotFound {
		c.JSON(http.StatusNotFound, controller.RenderError(errs.ErrMessageNotFound))
		return
	}
	if err != nil {
		controller.InternalServiceError(c, err)
		return
	}
	messageBodies := make([]messageBody, 0, len(messages))
	for _, value := range messages {
		messageBodies = append(messageBodies, messageBody{
			MessageID: value.MessageID,
			To:        value.To,
			Status:    value.Status,
			Body:      value.Body,
			At:        value.UpdatedAt.Unix(),
		})
	}
	c.JSON(http.StatusOK, messageHistoryResponse{
		Messages: messageBodies,
	})
}
