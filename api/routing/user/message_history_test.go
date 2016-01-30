package user

import (
	"encoding/json"
	"fmt"
	"github.com/danfang/portal-server/model"
	"github.com/danfang/portal-server/model/types"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var getMessagesDB gorm.DB

func init() {
	gin.SetMode(gin.TestMode)
	getMessagesDB, _ = gorm.Open("sqlite3", ":memory:")
	getMessagesDB.CreateTable(&model.User{}, &model.Message{})
}

func TestGetMessagesEndpoint_NoMessages(t *testing.T) {
	w := testGetMessages(404)
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"messages":[]}`, w.Body.String())
}

func TestGetMessagesEndpoint_AllMessages(t *testing.T) {
	user := model.User{
		Email: "test@portal.com",
	}
	getMessagesDB.Create(&user)
	message := &model.Message{
		User:      user,
		To:        "justin",
		MessageID: "1",
		Body:      "hello",
		Status:    types.MessageStatusDelivered.String(),
	}
	getMessagesDB.Create(&message)
	w := testGetMessages(user.ID)
	assert.Equal(t, 200, w.Code)
	var res messageHistoryResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, 1, len(res.Messages))
	assert.Equal(t, "justin", res.Messages[0].To)
	assert.Equal(t, "1", res.Messages[0].MessageID)
	assert.Equal(t, "hello", res.Messages[0].Body)
	assert.Equal(t, "delivered", res.Messages[0].Status)
	assert.Equal(t, message.UpdatedAt.Unix(), res.Messages[0].At)
}

func TestGetMessagesEndpoint_AllMessagesLimit(t *testing.T) {
	user := model.User{
		Email: "test2@portal.com",
	}
	getMessagesDB.Create(&user)
	messagesCreated := 3 * messageHistoryLimit
	for i := 1; i <= messagesCreated; i++ {
		getMessagesDB.Create(&model.Message{
			User:      user,
			To:        "myself",
			MessageID: fmt.Sprintf("message%d", i),
			Body:      "goodbye",
			Status:    types.MessageStatusDelivered.String(),
		})
	}
	w := testGetMessages(user.ID)
	assert.Equal(t, 200, w.Code)
	var res messageHistoryResponse
	assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
	assert.Equal(t, messageHistoryLimit, len(res.Messages))
	// Make sure they are in chronological order
	for i := 0; i < messageHistoryLimit; i++ {
		expectedMid := fmt.Sprintf("message%d", messagesCreated-i)
		assert.Equal(t, expectedMid, res.Messages[i].MessageID)
		assert.Equal(t, "myself", res.Messages[i].To)
		assert.Equal(t, "goodbye", res.Messages[i].Body)
		assert.Equal(t, "delivered", res.Messages[0].Status)
	}
}

func testGetMessages(userID uint) *httptest.ResponseRecorder {
	// Create the router
	userRouter := Router{&getMessagesDB, http.DefaultClient}
	r := gin.New()

	// Set the userID
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	r.GET("/", userRouter.GetMessageHistoryEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}
