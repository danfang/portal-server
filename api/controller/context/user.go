package context

import (
	"github.com/gin-gonic/gin"
	"portal-server/model"
)

const userKey = "user"

// UserToContext injects a user into the context <userKey, user>
func UserToContext(c *gin.Context, user *model.User) {
	c.Set(userKey, user)
}

// UserFromContext retrieves a user from the current context
func UserFromContext(c *gin.Context) *model.User {
	return c.MustGet(userKey).(*model.User)
}
