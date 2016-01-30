package user

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

// Router handles authenticated user operations.
type Router struct {
	Db         *gorm.DB
	HTTPClient *http.Client
}
