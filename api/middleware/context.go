package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portal-server/api/controller/context"
	"portal-server/store"
)

func SetStore(s store.Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		context.StoreToContext(c, s)
		c.Next()
	}
}

func SetWebClient(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		context.WebClientToContext(c, client)
		c.Next()
	}
}
