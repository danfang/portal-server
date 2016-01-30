package user

import (
	"github.com/jinzhu/gorm"
	"net/http"
)

// Router handles authenticated user operations.
type Router struct {
	Db         *gorm.DB
	HTTPClient *http.Client
}
