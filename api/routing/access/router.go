package access

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

// Router handles access routes for user creation and authentication.
type Router struct {
	Db         *gorm.DB
	HTTPClient *http.Client
}
