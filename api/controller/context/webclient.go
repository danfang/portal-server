package context

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portal-server/api/util"
)

const wcKey = "webclient"

func WebClientToContext(c *gin.Context, client *http.Client) {
	c.Set(wcKey, client)
}

func WebClientFromContext(c *gin.Context, url string) *util.WebClient {
	return &util.WebClient{
		BaseURL:    url,
		HTTPClient: c.MustGet(wcKey).(*http.Client),
	}
}
