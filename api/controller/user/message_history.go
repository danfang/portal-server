package user

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/model"

	"github.com/gin-gonic/gin"
)

const messageHistoryLimit = 1000

type messageHistoryResponse struct {
	Messages []messageBody `json:"messages"`
}

type messageBody struct {
	MessageID string `json:"mid"`
	To        string `json:"to"`
	Status    string `json:"status"`
	Body      string `json:"body"`
	At        int64  `json:"at"`
}

// GetMessageHistoryEndpoint retrieves user messages up to a
// given limit.
func (r Router) GetMessageHistoryEndpoint(c *gin.Context) {
	user := c.MustGet("user").(*model.User)
	messages, err := r.Store.Messages().GetMessagesByUser(user, messageHistoryLimit)
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
