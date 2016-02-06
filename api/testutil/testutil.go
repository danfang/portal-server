package testutil

import (
	"github.com/gin-gonic/gin"
)

func TestRouter(middleware ...gin.HandlerFunc) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware...)
	return r
}
