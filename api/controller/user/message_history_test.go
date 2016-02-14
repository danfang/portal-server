package user

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"portal-server/api/controller/context"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/model"
	"portal-server/store"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMessageHistory(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)
	g.Describe("GET /user/message/history", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownTestStore(s)
		})

		g.It("Should return an empty array for a user with no messages", func() {
			w := testGetMessages(s, &model.User{})
			assert.Equal(t, 200, w.Code)
			assert.JSONEq(t, `{"messages":[]}`, w.Body.String())
		})

		g.It("Should return all messages for a user", func() {
			user := model.User{Email: "test@portal.com"}
			s.Users().CreateUser(&user)
			message := &model.Message{
				User:      user,
				To:        "justin",
				MessageID: "1",
				Body:      "hello",
				Status:    model.MessageStatusDelivered,
			}
			s.Messages().CreateMessage(message)
			w := testGetMessages(s, &user)
			assert.Equal(t, 200, w.Code)
			var res messageHistoryResponse
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Equal(t, 1, len(res.Messages))
			assert.Equal(t, "justin", res.Messages[0].To)
			assert.Equal(t, "1", res.Messages[0].MessageID)
			assert.Equal(t, "hello", res.Messages[0].Body)
			assert.Equal(t, "delivered", res.Messages[0].Status)
			assert.Equal(t, message.UpdatedAt.Unix(), res.Messages[0].At)
		})

		g.It("Should return at most messageHistoryLimit messages", func() {
			user := model.User{Email: "test@portal.com"}
			s.Users().CreateUser(&user)
			messagesCreated := 3 * messageHistoryLimit
			for i := 1; i <= messagesCreated; i++ {
				s.Messages().CreateMessage(&model.Message{
					User:      user,
					To:        "myself",
					MessageID: fmt.Sprintf("message%d", i),
					Body:      "goodbye",
					Status:    model.MessageStatusDelivered,
				})
			}
			w := testGetMessages(s, &user)
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
		})
	})
}

func testGetMessages(s store.Store, user *model.User) *httptest.ResponseRecorder {
	r := testutil.TestRouter(middleware.SetStore(s))

	// Set the user context
	r.Use(func(c *gin.Context) {
		context.UserToContext(c, user)
		c.Next()
	})

	r.GET("/", GetMessageHistoryEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}
