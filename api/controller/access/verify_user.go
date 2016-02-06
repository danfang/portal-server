package access

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/errs"
	"portal-server/model"
	"time"

	"github.com/gin-gonic/gin"
	"portal-server/api/controller/context"
	"portal-server/store"
)

type verificationToken struct {
	Token string `json:"token"`
}

// VerifyUserEndpoint handles a GET request that consumes a user's verification token
// for users who registered with an email and password.
func VerifyUserEndpoint(c *gin.Context) {
	s := context.StoreFromContext(c)
	s.Transaction(func(store store.Store) error {
		user, err := checkVerificationToken(store, c.Param("token"))
		if err != nil {
			c.JSON(http.StatusBadRequest, controller.RenderError(err))
			return nil
		}
		user.Verified = true
		if err := store.Users().SaveUser(user); err != nil {
			controller.InternalServiceError(c, err)
			return err
		}
		return nil
	})
	c.JSON(http.StatusOK, controller.RenderSuccess())
}

func checkVerificationToken(store store.Store, param string) (*model.User, error) {
	token, found := store.VerificationTokens().FindToken(&model.VerificationToken{
		Token: param,
	})

	if !found {
		return nil, errs.ErrInvalidVerificationToken
	}

	// Expired token
	if time.Now().After(token.ExpiresAt) {
		store.VerificationTokens().DeleteToken(token)
		return nil, errs.ErrExpiredVerificationToken
	}

	// Check for existing user account
	user, err := store.VerificationTokens().GetRelatedUser(token)
	if err != nil {
		store.VerificationTokens().DeleteToken(token)
		return nil, errs.ErrInvalidVerificationToken
	}

	store.VerificationTokens().DeleteToken(token)
	return user, nil
}
