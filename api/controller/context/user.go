package context

import (
	"github.com/gin-gonic/gin"
	"portal-server/model"
)

const userKey = "user"

func UserToContext(c *gin.Context, user *model.User) {
	c.Set(userKey, user)
}

func UserFromContext(c *gin.Context) *model.User {
	return c.MustGet(userKey).(*model.User)
}
