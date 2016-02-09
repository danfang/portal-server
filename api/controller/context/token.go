package context

import (
	"portal-server/model"

	"github.com/gin-gonic/gin"
)

const userTokenKey = "userToken"

// UserTokenToContext injects a user into the context <userKey, user>
func UserTokenToContext(c *gin.Context, user *model.UserToken) {
	c.Set(userTokenKey, user)
}

// UserTokenFromContext retrieves a user from the current context
func UserTokenFromContext(c *gin.Context) *model.UserToken {
	return c.MustGet(userTokenKey).(*model.UserToken)
}
