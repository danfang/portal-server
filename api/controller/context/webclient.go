package context

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portal-server/api/util"
)

const wcKey = "webclient"

// WebClientToContext sets the key value <wcKey, client>
func WebClientToContext(c *gin.Context, client *http.Client) {
	c.Set(wcKey, client)
}

// WebClientFromContext retrieves a WebClient given a url string from the context
func WebClientFromContext(c *gin.Context, url string) *util.WebClient {
	return &util.WebClient{
		BaseURL:    url,
		HTTPClient: c.MustGet(wcKey).(*http.Client),
	}
}
