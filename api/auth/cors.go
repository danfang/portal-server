package auth

import (
	"github.com/gin-gonic/gin"
	"strings"
)

// CORSMiddleware is Gin middleware that allows for fine-grained
// Cross-Origin control.
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		allowedHeaders := []string{
			"Content-Type",
			"Content-Length",
			"Cache-Control",
			"accept",
			"origin",
			UserTokenHeader,
			UserIDHeader,
		}
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}

		c.Next()
	}
}
