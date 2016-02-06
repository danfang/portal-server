package access

import (
	"net/http"
	"net/http/httptest"
	"portal-server/api/errs"
	"portal-server/model"
	"testing"
	"time"

	"github.com/franela/goblin"
	"github.com/stretchr/testify/assert"
	"portal-server/api/middleware"
	"portal-server/api/testutil"
	"portal-server/store"
)

var verifyUserStore = store.GetTestStore()
var verifyUser *model.User

func init() {
	verifyUser = &model.User{Email: "test@portal.com"}
	verifyUserStore.Users().CreateUser(verifyUser)
}

func TestVerifyUserEndpoint_NoToken(t *testing.T) {
	w := testVerifyUser("")
	assert.Equal(t, 404, w.Code)
}

func TestVerifyUser(t *testing.T) {
	var s store.Store
	g := goblin.Goblin(t)
	g.Describe("GET /verify/user", func() {
		g.BeforeEach(func() {
			s = store.GetTestStore()
		})

		g.AfterEach(func() {
			store.TeardownStoreForTest(s)
		})
	})
}

func TestVerifyUserEndpoint_BadToken(t *testing.T) {
	user := &model.User{
		Email: "test_endpoint_bad_token@test.com",
	}
	verifyUserStore.Users().CreateUser(user)
	token := &model.VerificationToken{
		User:      *user,
		Token:     "test_endpoint_bad_token",
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	verifyUserStore.VerificationTokens().CreateToken(token)
	w := testVerifyUser("invalid_token")
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrInvalidVerificationToken.Error())

	// Make sure token is not deleted
	tokenCount := verifyUserStore.VerificationTokens().GetCount(&model.VerificationToken{
		Token: "test_endpoint_bad_token",
	})
	assert.Equal(t, 1, tokenCount)
}

func TestVerifyUserEndpoint_ExpiredToken(t *testing.T) {
	user := &model.User{
		Email: "test_endpoint_expired_token@test.com",
	}
	verifyUserStore.Users().CreateUser(user)
	token := model.VerificationToken{
		User:      *user,
		Token:     "test_endpoint_expired_token",
		ExpiresAt: time.Now(),
	}
	verifyUserStore.VerificationTokens().CreateToken(&token)
	w := testVerifyUser("test_endpoint_expired_token")
	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), errs.ErrExpiredVerificationToken.Error())

	// Make sure token is not deleted
	tokenCount := verifyUserStore.VerificationTokens().GetCount(&model.VerificationToken{
		Token: "test_endpoint_expired_token",
	})
	assert.Equal(t, 0, tokenCount)
}

func TestVerifyUserEndpoint_ValidToken(t *testing.T) {
	user := &model.User{
		Email: "test_endpoint_valid_token@test.com",
	}
	verifyUserStore.Users().CreateUser(user)
	token := &model.VerificationToken{
		User:      *user,
		Token:     "test_endpoint_valid_token",
		ExpiresAt: time.Now().Add(1 * time.Minute),
	}
	verifyUserStore.VerificationTokens().CreateToken(token)
	w := testVerifyUser("test_endpoint_valid_token")
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `{"success":true}`, w.Body.String())

	// Make sure token is not deleted
	tokenCount := verifyUserStore.VerificationTokens().GetCount(&model.VerificationToken{
		Token: "test_endpoint_valid_token",
	})
	assert.Equal(t, 0, tokenCount)
}

func TestCheckVerificationToken_NoSuchToken(t *testing.T) {
	_, err := checkVerificationToken(verifyUserStore, "no_such_token")
	assert.EqualError(t, err, "invalid_verification_token")
}

func TestCheckVerificationToken_Expired(t *testing.T) {
	expiredToken := &model.VerificationToken{
		User:      *verifyUser,
		ExpiresAt: time.Now().Add(-1 * time.Minute),
		Token:     "expired_token",
	}
	verifyUserStore.VerificationTokens().CreateToken(expiredToken)

	_, err := checkVerificationToken(verifyUserStore, "expired_token")
	assert.EqualError(t, err, "expired_verification_token")

	deletedToken, _ := verifyUserStore.VerificationTokens().FindDeletedToken(&model.VerificationToken{
		Token: "expired_token",
	})
	assert.NotNil(t, deletedToken.DeletedAt)
}

func TestCheckVerificationToken_NoUserToken(t *testing.T) {
	noUserToken := model.VerificationToken{
		UserID:    uint(404),
		ExpiresAt: time.Now().Add(time.Minute),
		Token:     "no_user_token",
	}
	verifyUserStore.VerificationTokens().CreateToken(&noUserToken)

	_, err := checkVerificationToken(verifyUserStore, "no_user_token")
	assert.EqualError(t, err, "invalid_verification_token")

	deletedToken, _ := verifyUserStore.VerificationTokens().FindDeletedToken(&model.VerificationToken{
		Token: "no_user_token",
	})
	assert.NotNil(t, deletedToken.DeletedAt)
}

func TestCheckVerificationToken_ValidToken(t *testing.T) {
	token := model.VerificationToken{
		User:      *verifyUser,
		ExpiresAt: time.Now().Add(time.Second),
		Token:     "valid_token",
	}
	verifyUserStore.VerificationTokens().CreateToken(&token)

	fromDB, err := checkVerificationToken(verifyUserStore, "valid_token")
	assert.NoError(t, err)
	assert.Equal(t, verifyUser.ID, fromDB.ID)

	deletedToken, _ := verifyUserStore.VerificationTokens().FindDeletedToken(&model.VerificationToken{
		Token: "valid_token",
	})
	assert.NotNil(t, deletedToken.DeletedAt)
}

func testVerifyUser(token string) *httptest.ResponseRecorder {
	// Create the router
	r := testutil.TestRouter(middleware.SetStore(verifyUserStore))

	// Test the response
	r.GET("/:token", VerifyUserEndpoint)
	w := httptest.NewRecorder()

	// Send the input
	req, _ := http.NewRequest("GET", "/"+token, nil)
	r.ServeHTTP(w, req)
	return w
}
