package user

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"portal-server/store"
)

// Router handles authenticated user operations.
type Router struct {
	Store      store.Store
	HTTPClient *http.Client
}
