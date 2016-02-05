package user

import (
	"net/http"
	"portal-server/api/routing"
	"portal-server/model"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
	userID := c.MustGet("userID").(uint)
	messages, err := getMessages(r.Db, userID)
	if err != nil {
		routing.InternalServiceError(c, err)
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

func getMessages(db *gorm.DB, userID uint) ([]model.Message, error) {
	var messages []model.Message
	if err := db.Where(&model.Message{UserID: userID}).Order("id desc").
		Limit(messageHistoryLimit).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
