package middleware

import (
	"net/http"
	"portal-server/api/controller"
	"portal-server/api/controller/context"
	"portal-server/api/errs"
	"portal-server/model"
	"portal-server/store"
	"time"

	"github.com/gin-gonic/gin"
)

// Headers for user authentication
const (
	UserTokenHeader = "X-USER-TOKEN"
	UserIDHeader    = "X-USER-ID"
)

// AuthenticationMiddleware handles authentication for protected user
// endpoints by checking for valid user id and user token headers.
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for valid headers
		token := c.Request.Header.Get(UserTokenHeader)
		userUUID := c.Request.Header.Get(UserIDHeader)
		if token == "" || userUUID == "" {
			c.JSON(http.StatusUnauthorized, controller.RenderError(errs.ErrMissingHeaders))
			c.Abort()
			return
		}

		// Check for valid token for the given user
		store := context.StoreFromContext(c)
		user, userToken, err := authenticate(store, token, userUUID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, controller.RenderError(err))
			c.Abort()
			return
		}

		context.UserToContext(c, user)
		context.UserTokenToContext(c, userToken)
		c.Next()
	}
}

func authenticate(store store.Store, token, uuid string) (*model.User, *model.UserToken, error) {
	// User not found
	user, found := store.Users().FindUser(&model.User{UUID: uuid})
	if !found {
		return nil, nil, errs.ErrInvalidUserToken
	}

	// Token not found
	userToken, found := store.UserTokens().FindToken(&model.UserToken{Token: token, UserID: user.ID})
	if !found {
		return nil, nil, errs.ErrInvalidUserToken
	}

	// Token expired
	if !userToken.ExpiresAt.IsZero() && time.Now().After(userToken.ExpiresAt) {
		store.UserTokens().DeleteToken(userToken)
		return nil, nil, errs.ErrInvalidUserToken
	}

	// Account not verified
	if !user.Verified {
		return nil, nil, errs.ErrAccountNotVerified
	}

	return user, userToken, nil
}
