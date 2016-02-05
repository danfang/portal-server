package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var corsRouter *gin.Engine

func init() {
	gin.SetMode(gin.TestMode)
	corsRouter = gin.New()
	corsRouter.Use(CORSMiddleware())
	corsRouter.GET("/", func(c *gin.Context) {
		c.String(200, "done")
	})
}

func checkHeaders(t *testing.T, header http.Header) {
	assert.Equal(t, "*", header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", header.Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "POST, GET, OPTIONS, PUT, DELETE", header.Get("Access-Control-Allow-Methods"))

	allowedHeaders := header.Get("Access-Control-Allow-Headers")
	assert.Contains(t, allowedHeaders, "Content-Type")
	assert.Contains(t, allowedHeaders, "Cache-Control")
	assert.Contains(t, allowedHeaders, "Content-Length")
	assert.Contains(t, allowedHeaders, "accept")
	assert.Contains(t, allowedHeaders, "origin")
	assert.Contains(t, allowedHeaders, "X-USER-ID")
	assert.Contains(t, allowedHeaders, "X-USER-TOKEN")
}

func TestCORS_Headers(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	corsRouter.ServeHTTP(w, req)
	checkHeaders(t, w.Header())
	assert.Equal(t, "done", w.Body.String())
}

func TestCORS_Options(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/", nil)
	w := httptest.NewRecorder()
	corsRouter.ServeHTTP(w, req)
	checkHeaders(t, w.Header())
	assert.Empty(t, w.Body.String())
}
