package middleware

import (
	"net/http"
	"portal-server/api/controller/context"
	"portal-server/store"

	"github.com/gin-gonic/gin"
)

// SetStore injects the datastore into every gin context
func SetStore(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		context.StoreToContext(c, s)
		c.Next()
	}
}

// SetWebClient injects an HTTP client into every gin context
func SetWebClient(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		context.WebClientToContext(c, client)
		c.Next()
	}
}
