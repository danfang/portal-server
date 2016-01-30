package auth

import (
	"github.com/danfang/portal-server/api/errs"
	"github.com/danfang/portal-server/api/routing"
	"github.com/danfang/portal-server/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

// Headers for user authentication
const (
	UserTokenHeader = "X-USER-TOKEN"
	UserIDHeader    = "X-USER-ID"
)

// AuthenticationMiddleware handles authentication for protected user
// endpoints by checking for valid user id and user token headers.
func AuthenticationMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for valid headers
		token := c.Request.Header.Get(UserTokenHeader)
		userUUID := c.Request.Header.Get(UserIDHeader)
		if token == "" || userUUID == "" {
			c.JSON(http.StatusUnauthorized, routing.RenderError(errs.ErrMissingHeaders))
			c.Abort()
			return
		}

		// Check for valid token for the given user
		userID, err := authenticate(db, token, userUUID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, routing.RenderError(err))
			c.Abort()
			return
		}

		c.Set("userID", userID)
		c.Next()
	}
}

func authenticate(db *gorm.DB, token, userUUID string) (uint, error) {
	var userToken model.UserToken
	var user model.User
	// User not found
	if db.Where(&model.User{UUID: userUUID}).First(&user).RecordNotFound() {
		return 0, errs.ErrInvalidUserToken
	}

	// Token not found
	if db.Where(&model.UserToken{Token: token, UserID: user.ID}).First(&userToken).RecordNotFound() {
		return 0, errs.ErrInvalidUserToken
	}

	// Token expired
	if !userToken.ExpiresAt.IsZero() && time.Now().After(userToken.ExpiresAt) {
		db.Delete(&userToken)
		return 0, errs.ErrInvalidUserToken
	}

	// Acount not verified
	if !user.Verified {
		return 0, errs.ErrAccountNotVerified
	}
	return user.ID, nil
}
