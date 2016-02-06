package testutil

import (
	"github.com/gin-gonic/gin"
)

// TestRouter is a convenience method to construct a gin router from
// middleware, in test mode.
func TestRouter(middleware ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware...)
	return r
}
