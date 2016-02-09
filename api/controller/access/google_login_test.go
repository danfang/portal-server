package access

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"portal-server/api/errs"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/api/util"
	"portal-server/model"
	"portal-server/store"
	"testing"

	"github.com/franela/goblin"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

func TestGoogleLogin(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)

	g.Describe("POST /login/google", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownStoreForTest(s)
		})

		g.It("Should return 400 on invalid JSON input", func() {
			input := map[string]string{"id_token": ""}
			w := testGoogleLogin(s, input, 200, "")
			assert.Equal(t, 400, w.Code)
			assert.Contains(t, w.Body.String(), errs.ErrInvalidJSON.Error())
		})

		g.It("Should return 400 on invalid ID token", func() {
			input := map[string]string{"id_token": "token"}
			w := testGoogleLogin(s, input, 400, "{}")
			assert.Equal(t, 400, w.Code)
			assert.Contains(t, w.Body.String(), errs.ErrInvalidGoogleIDToken.Error())
		})

		g.It("Should return 500 on Google 404 response", func() {
			input := map[string]string{"id_token": "token"}
			w := testGoogleLogin(s, input, 404, "")
			assert.Equal(t, 500, w.Code)
			assert.Contains(t, w.Body.String(), errs.ErrGoogleOAuthUnavailable.Error())
		})

		g.It("Should return 400 and error if user Google account is unverified", func() {
			input := map[string]string{"id_token": "token"}
			output := util.GoogleUser{
				Sub:           "1000",
				Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
				Email:         "test@google.com",
				EmailVerified: "false",
			}
			w := testGoogleLogin(s, input, 200, output)
			assert.Equal(t, 400, w.Code)
			assert.Contains(t, w.Body.String(), errs.ErrGoogleAccountNotVerified.Error())
		})

		g.It("Should create a new user and linked account on valid credentials", func() {
			input := map[string]string{"id_token": "token"}
			output := util.GoogleUser{
				Sub:           "valid_user_sub",
				Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
				Email:         "test@google.com",
				EmailVerified: "true",
			}
			w := testGoogleLogin(s, input, 200, output)
			// Check login response
			assert.Equal(t, 200, w.Code)
			assertValidLoginResponse(t, w)

			// Check linked account is in DB
			linkedAccount, _ := s.LinkedAccounts().FindAccount(&model.LinkedAccount{
				AccountID: "valid_user_sub",
				Type:      model.LinkedAccountTypeGoogle,
			})

			assert.Equal(t, "valid_user_sub", linkedAccount.AccountID)

			// Check user is created
			user, _ := s.LinkedAccounts().GetRelatedUser(linkedAccount)
			assert.Equal(t, "test@google.com", user.Email)
			assert.True(t, user.Verified)
		})

		g.It("Should only create a new linked account if the Google email matches a user's", func() {
			user := &model.User{
				UUID:     uuid.NewV4().String(),
				Email:    "test2@google.com",
				Verified: false,
				Password: "my_password_hash",
			}
			s.Users().CreateUser(user)
			input := map[string]string{"id_token": "token"}
			output := util.GoogleUser{
				Sub:           "existing_user_sub",
				Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
				Email:         "test2@google.com",
				EmailVerified: "true",
			}
			w := testGoogleLogin(s, input, 200, output)
			// Check login response
			assert.Equal(t, 200, w.Code)
			assertValidLoginResponse(t, w)

			// Check linked account is in DB
			linkedAccount, _ := s.LinkedAccounts().FindAccount(&model.LinkedAccount{
				AccountID: "existing_user_sub",
				Type:      model.LinkedAccountTypeGoogle,
			})

			assert.Equal(t, "existing_user_sub", linkedAccount.AccountID)

			// Check user is created
			fromDB, _ := s.LinkedAccounts().GetRelatedUser(linkedAccount)
			assert.Equal(t, "test2@google.com", fromDB.Email)
			assert.True(t, fromDB.Verified)

			// Check that password login is disabled
			assert.Equal(t, "", fromDB.Password)
		})

		g.It("Should login without creating new accounts for existing users", func() {
			user := model.User{
				UUID:     uuid.NewV4().String(),
				Email:    "test3@google.com",
				Password: "some_password",
			}
			s.Users().CreateUser(&user)
			account := model.LinkedAccount{
				User:      user,
				AccountID: "existing_user_and_account_sub",
				Type:      model.LinkedAccountTypeGoogle,
			}
			s.LinkedAccounts().CreateAccount(&account)
			input := map[string]string{
				"id_token": "token",
			}
			output := util.GoogleUser{
				Sub:           "existing_user_and_account_sub",
				Aud:           "1045304436932-9vtokstg18sq2hu26hipueithq7sb0bq.apps.googleusercontent.com",
				Email:         "test3@google.com",
				EmailVerified: "true",
			}
			w := testGoogleLogin(s, input, 200, output)
			// Check login response
			assert.Equal(t, 200, w.Code)
			assertValidLoginResponse(t, w)

			// Check linked account is in DB
			linkedAccountCount := s.LinkedAccounts().GetCount(&model.LinkedAccount{
				AccountID: "existing_user_and_account_sub",
				Type:      model.LinkedAccountTypeGoogle,
			})

			assert.Equal(t, 1, linkedAccountCount)
		})
	})

	g.Describe("Google login data manipulation", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownStoreForTest(s)
		})

		g.It("Should successfully create a linked account given a Google user", func() {
			googleUser := util.GoogleUser{
				GivenName:     "Jon",
				FamilyName:    "Snow",
				EmailVerified: "true",
				Email:         "google@google.com",
				Sub:           "10000",
			}
			user, err := createLinkedGoogleAccount(s, &googleUser)
			assert.NoError(t, err)

			fromDB, _ := s.Users().FindUser(&model.User{Email: "google@google.com"})
			assert.Equal(t, fromDB.FirstName, "Jon")
			assert.Equal(t, fromDB.LastName, "Snow")
			assert.True(t, fromDB.Verified)

			var linkedAccount model.LinkedAccount
			s.Users().GetRelated(user, &linkedAccount)
			assert.Equal(t, linkedAccount.AccountID, "10000")
			assert.Equal(t, linkedAccount.Type, "google")
		})

		g.It("Should only create a linked account for an existing user", func() {
			original := model.User{
				Email:     "stannis@portal.com",
				FirstName: "Stannis",
				LastName:  "Baratheon",
				Verified:  false,
				Password:  "my_password",
			}

			s.Users().CreateUser(&original)

			googleUser := util.GoogleUser{
				GivenName:     "Stan",
				FamilyName:    "The Mannis",
				EmailVerified: "true",
				Email:         "stannis@portal.com",
				Sub:           "12345",
			}

			user, err := createLinkedGoogleAccount(s, &googleUser)
			assert.NoError(t, err)

			fromDB, _ := s.Users().FindUser(&model.User{Email: "stannis@portal.com"})

			assert.Equal(t, "Stannis", fromDB.FirstName)
			assert.Equal(t, "Baratheon", fromDB.LastName)

			// Check that the account is now verified and password disabled.
			assert.True(t, fromDB.Verified)
			assert.Equal(t, "", fromDB.Password)

			var linkedAccount model.LinkedAccount
			s.Users().GetRelated(user, &linkedAccount)
			assert.Equal(t, "12345", linkedAccount.AccountID)
			assert.Equal(t, model.LinkedAccountTypeGoogle, linkedAccount.Type)
		})

		g.It("Should only retrieve a user if the user already exists", func() {
			googleAccountID := "10101"

			original := model.User{
				Email: "existing@portal.com",
			}

			s.Users().CreateUser(&original)

			linkedAccount := model.LinkedAccount{
				User:      original,
				AccountID: googleAccountID,
				Type:      model.LinkedAccountTypeGoogle,
			}

			s.LinkedAccounts().CreateAccount(&linkedAccount)

			googleUser := util.GoogleUser{
				Sub:   googleAccountID,
				Email: "otherEmail@otherDomain.com",
			}

			// Make sure no data is modified
			user, err := createLinkedGoogleAccount(s, &googleUser)
			assert.NoError(t, err)
			assert.Equal(t, original.ID, user.ID)
			assert.Equal(t, original.Email, user.Email)

			// Make sure no new linked account is created.
			count := s.LinkedAccounts().GetCount(&model.LinkedAccount{
				AccountID: googleAccountID,
				Type:      model.LinkedAccountTypeGoogle,
			})
			assert.Equal(t, 1, count)
		})

	})
}

func testGoogleLogin(s store.Store, input interface{}, code int, response interface{}) *httptest.ResponseRecorder {
	// Setup mock Google server/client
	output, _ := json.Marshal(response)
	server, client := util.TestHTTP(func(*http.Request) {}, code, string(output))
	googleOAuthEndpoint = server.URL

	// Setup router
	r := testutil.TestRouter(
		middleware.SetWebClient(client.HTTPClient),
		middleware.SetStore(s),
	)

	// Setup endpoints
	r.POST("/", GoogleLoginEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(string(body)))
	r.ServeHTTP(w, req)
	return w
}
