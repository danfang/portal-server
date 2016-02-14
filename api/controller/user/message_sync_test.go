package user

import (
	"net/http"
	"net/http/httptest"
	"portal-server/api/controller/context"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/model"
	"portal-server/store"
	"testing"

	"encoding/json"
	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMessageSync(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)
	g.Describe("GET /user/message/sync/:mid", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownStoreForTest(s)
		})

		g.It("Should give an error for a non-existent message id", func() {
			w := testSyncMessages(s, &model.User{}, "bad_mid")
			assert.Equal(t, 400, w.Code)
		})

		g.It("Should return all messages for a user since the given message id", func() {
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
			message2 := &model.Message{
				User:      user,
				To:        "justin",
				MessageID: "2",
				Body:      "hello",
				Status:    model.MessageStatusDelivered,
			}
			s.Messages().CreateMessage(message2)
			message3 := &model.Message{
				User:      user,
				To:        "justin",
				MessageID: "3",
				Body:      "hello",
				Status:    model.MessageStatusDelivered,
			}
			s.Messages().CreateMessage(message3)
			message4 := &model.Message{
				User:      user,
				To:        "justin",
				MessageID: "4",
				Body:      "hello",
				Status:    model.MessageStatusDelivered,
			}
			s.Messages().CreateMessage(message4)
			w := testSyncMessages(s, &user, message3.MessageID)
			assert.Equal(t, 200, w.Code)
			var res messageHistoryResponse
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Equal(t, 1, len(res.Messages))
			assert.Equal(t, "justin", res.Messages[0].To)
			assert.Equal(t, message4.MessageID, res.Messages[0].MessageID)
			assert.Equal(t, "hello", res.Messages[0].Body)
			assert.Equal(t, "delivered", res.Messages[0].Status)
			assert.Equal(t, message.UpdatedAt.Unix(), res.Messages[0].At)
		})
	})
}

func testSyncMessages(s store.Store, user *model.User, messageID string) *httptest.ResponseRecorder {
	r := testutil.TestRouter(middleware.SetStore(s))

	// Set the user context
	r.Use(func(c *gin.Context) {
		context.UserToContext(c, user)
		c.Next()
	})

	r.GET("/:mid", SyncMessagesEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/"+messageID, nil)
	r.ServeHTTP(w, req)
	return w
}
