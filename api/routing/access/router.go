package access

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

// Router handles access routes for user creation and authentication.
type Router struct {
	Db         *gorm.DB
	HTTPClient *http.Client
}
