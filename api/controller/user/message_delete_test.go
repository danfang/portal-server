package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/api/errs"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/model"
	"portal-server/store"
	"testing"

	"github.com/franela/goblin"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestMessageDelete(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)

	g.Describe("DELETE /user/message/:mid", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownTestStore(s)
		})

		g.It("Should give a 404 for a non-existent message id", func() {
			w := testDeleteMessage(s, &model.User{}, "bad_mid")
			assert.Equal(t, 404, w.Code)

			var res controller.Error
			json.Unmarshal(w.Body.Bytes(), &res)
			assert.Equal(t, errs.ErrMessageNotFound.Error(), res.Error)
		})

		g.It("Should delete the given message id", func() {
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

			w := testDeleteMessage(s, &user, message2.MessageID)
			assert.Equal(t, 200, w.Code)

			var res controller.SuccessResponse
			assert.NoError(t, json.Unmarshal(w.Body.Bytes(), &res))
			assert.Equal(t, true, res.Success)

			messages, _ := s.Messages().GetMessagesByUser(&user, 10)
			assert.Equal(t, 1, len(messages))
			assert.Equal(t, "1", messages[0].MessageID)
		})
	})
}

func testDeleteMessage(s store.Store, user *model.User, messageID string) *httptest.ResponseRecorder {
	r := testutil.TestRouter(middleware.SetStore(s))

	// Set the user context
	r.Use(func(c *gin.Context) {
		context.UserToContext(c, user)
		c.Next()
	})

	r.DELETE("/:mid", DeleteMessageEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("DELETE", "/"+messageID, nil)
	r.ServeHTTP(w, req)
	return w
}
