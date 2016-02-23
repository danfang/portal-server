package user

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestAddContacts(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)

	g.Describe("POST /user/contacts", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownTestStore(s)
		})

		g.It("Should successfully upload contacts", func() {
			user := model.User{
				UUID:  "1",
				Email: "hello@world.com",
			}
			s.Users().CreateUser(&user)
			numContacts := 5

			contactsJson := make([]map[string]interface{}, 0, numContacts)
			for i := 0; i < numContacts; i++ {
				contact := map[string]interface{}{
					"name": fmt.Sprintf("contact%d", i),
					"cid":  uuid.NewV4().String(),
					"phone_numbers": []map[string]string{
						{
							"name":   fmt.Sprintf("home%d", i),
							"number": fmt.Sprintf("homenumber%d", i),
						},
						{
							"name":   fmt.Sprintf("cell%d", i),
							"number": fmt.Sprintf("cellnumber%d", i),
						},
					},
				}
				contactsJson = append(contactsJson, contact)

			}
			input := map[string]interface{}{
				"contacts": contactsJson,
			}
			w := testAddContacts(s, &user, input)
			assert.Equal(t, 200, w.Code)

			var response controller.SuccessResponse
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.True(t, response.Success)

			for i := 0; i < numContacts; i++ {
				contact, _ := s.Contacts().FindContact(&model.Contact{
					UserID: user.ID,
					Name:   fmt.Sprintf("contact%d", i),
				})
				assert.Equal(t, 2, len(contact.PhoneNumbers))
			}
		})

		g.It("Should return 400 on missing JSON fields", func() {
			user := model.User{
				UUID:  "1",
				Email: "hello@world.com",
			}
			s.Users().CreateUser(&user)
			input := map[string]interface{}{
				"contacts": map[string]interface{}{
					"name": "contact",
					"uuid": uuid.NewV4().String(),
					"phone_numbers": []map[string]string{
						{
							"name":   "home",
							"number": "homenumber",
						},
						{
							"name":   "cell",
							"number": "cell1number",
						},
					},
				},
			}
			w := testAddContacts(s, &user, input)
			assert.Equal(t, 400, w.Code)
			var err controller.Error
			json.Unmarshal(w.Body.Bytes(), &err)
			assert.Equal(t, errs.ErrInvalidJSON.Error(), err.Error)
		})

		g.It("Should return 400 on invalid JSON fields", func() {
			user := model.User{
				UUID:  "1",
				Email: "hello@world.com",
			}
			s.Users().CreateUser(&user)
			input := map[string]interface{}{
				"contacts": map[string]interface{}{
					"name": "contact",
					"cid":  "not-a-uuid",
					"phone_numbers": []map[string]string{
						{
							"name":   "home",
							"number": "homenumber",
						},
						{
							"name":   "cell",
							"number": "cell1number",
						},
					},
				},
			}
			w := testAddContacts(s, &user, input)
			assert.Equal(t, 400, w.Code)
			var err controller.Error
			json.Unmarshal(w.Body.Bytes(), &err)
			assert.Equal(t, errs.ErrInvalidJSON.Error(), err.Error)
		})
	})
}

func testAddContacts(s store.Store, user *model.User, input interface{}) *httptest.ResponseRecorder {
	r := testutil.TestRouter(middleware.SetStore(s))
	r.Use(func(c *gin.Context) {
		context.UserToContext(c, user)
		c.Next()
	})
	r.POST("/", AddContactsEndpoint)
	w := httptest.NewRecorder()
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(body)))
	r.ServeHTTP(w, req)
	return w
}
