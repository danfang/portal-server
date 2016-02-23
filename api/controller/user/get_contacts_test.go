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
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetContacts(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)

	g.Describe("GET /user/contacts", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownTestStore(s)
		})

		g.It("Should successfully get an empty array on no contacts", func() {
			user := model.User{
				UUID:  "1",
				Email: "hello@world.com",
			}
			s.Users().CreateUser(&user)
			w := testGetContacts(s, &user)
			assert.Equal(t, 200, w.Code)

			var response contactsJson
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, 0, len(response.Contacts))
		})

		g.It("Should successfully get contacts", func() {
			user := model.User{
				UUID:  "1",
				Email: "hello@world.com",
			}
			s.Users().CreateUser(&user)
			numContacts := 5
			for i := 0; i < numContacts; i++ {
				contact := model.Contact{
					UserID: user.ID,
					Name:   fmt.Sprintf("contact%d", i),
					UUID:   uuid.NewV4().String(),
					PhoneNumbers: []model.ContactPhone{
						{
							Name:   fmt.Sprintf("home%d", i),
							Number: fmt.Sprintf("homenumber%d", i),
						},
						{
							Name:   fmt.Sprintf("cell%d", i),
							Number: fmt.Sprintf("cellnumber%d", i),
						},
					},
				}
				s.Contacts().CreateContact(&contact)

			}
			w := testGetContacts(s, &user)
			assert.Equal(t, 200, w.Code)

			var response contactsJson
			json.Unmarshal(w.Body.Bytes(), &response)
			assert.Equal(t, 5, len(response.Contacts))
			for i := 0; i < numContacts; i++ {
				contact := response.Contacts[i]
				assert.Equal(t, 2, len(contact.PhoneNumbers))
			}
		})
	})
}

func testGetContacts(s store.Store, user *model.User) *httptest.ResponseRecorder {
	r := testutil.TestRouter(middleware.SetStore(s))
	r.Use(func(c *gin.Context) {
		context.UserToContext(c, user)
		c.Next()
	})
	r.GET("/", GetContactsEndpoint)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)
	return w
}
