package context

import (
	"github.com/gin-gonic/gin"
	"portal-server/store"
)

const storeKey = "store"

func StoreToContext(c *gin.Context, s store.Store) {
	c.Set(storeKey, s)
}

func StoreFromContext(c *gin.Context) store.Store {
	return c.MustGet(storeKey).(store.Store)
}
